#!/bin/bash

set -e

APP_NAME=project_template_go
configFIleName=config.yaml
APP_DIR="/${APP_NAME}"

SYSTEMD_PATH="/etc/systemd/system"

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
PARENT_DIR=$(dirname ${SCRIPT_DIR})

echo parent dir is:$PARENT_DIR

function clear_install_files() {
    if [ -d "${PARENT_DIR}" ]; then
        find "${PARENT_DIR}" -type f -name "${APP_NAME}*" -exec rm -f {} \;
        find "${PARENT_DIR}" -type d -name "${APP_NAME}*" -empty -delete
        rm -rf ${SCRIPT_DIR}
    fi
}

function install() {

    if [[ -f ${SCRIPT_DIR}/${APP_NAME} ]] && [[ -f ${SCRIPT_DIR}/${APP_NAME}.service ]] && [[ -f ${SCRIPT_DIR}/conf/${configFIleName} ]]; then
        if [ ! -d "${APP_DIR}" ]; then
            mkdir -p ${APP_DIR}
            mkdir -p ${APP_DIR}/logs
            mkdir -p ${APP_DIR}/conf
        fi

        # install binary
        TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
        if [[ -f "${APP_DIR}/${APP_NAME}" ]]; then
            mv "${APP_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}.bak_${TIMESTAMP}"
        fi
        cp -f "${SCRIPT_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}"

        # configure file
        if [[ -f "${APP_DIR}/conf/${configFIleName}" ]]; then
            mv "${APP_DIR}/conf/${configFIleName}" "${APP_DIR}/conf/${configFIleName}.bak_${TIMESTAMP}"
        fi
        cp -f ${SCRIPT_DIR}/conf/${configFIleName} ${APP_DIR}/conf/${configFIleName}

        # install systemd file to SYSTEMD_PATH and restart service
        service_file=${SCRIPT_DIR}/${APP_NAME}.service
        if [[ -f "${service_file}" ]]; then
            cp -f "${service_file}" "${SYSTEMD_PATH}"
            service_name=$(basename "${service_file}" .service)
            systemctl enable "${service_name}"
            systemctl daemon-reload
            systemctl start "${service_name}.service"
        fi
    else
        echo "installation file error"
        exit 1
    fi
}

function updateBinary() {
    read -p "sure want to update ${APP_NAME} [yes/no]：" flag
    if [ -z $flag ]; then
        echo "input error" && exit 1
    elif [ "$flag" = "yes" -o "$flag" = "ye" -o "$flag" = "y" ]; then
        # update binary
        TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
        systemctl stop "${APP_NAME}.service"
        if [[ -f "${APP_DIR}/${APP_NAME}" ]]; then
            mv "${APP_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}.bak_${TIMESTAMP}"
        fi
        cp -f "${SCRIPT_DIR}/${APP_NAME}" "${APP_DIR}/${APP_NAME}"

        # update systemd service
        service_file=${SCRIPT_DIR}/${APP_NAME}.service
        if [[ -f "${service_file}" ]]; then
            cp -f "${service_file}" "${SYSTEMD_PATH}"
            systemctl daemon-reload
            systemctl restart "${APP_NAME}.service"
        fi
    fi
}

function uninstall() {
    read -p "sure want to uninstall ${APP_NAME} [yes/no]：" flag
    if [ -z "$flag" ]; then
        echo "input error" && exit 1
    elif [ "$flag" = "yes" ] || [ "$flag" = "ye" ] || [ "$flag" = "y" ]; then
        for service_file in ${SYSTEMD_PATH}/${APP_NAME}*.service; do
            if [[ -f ${service_file} ]]; then
                service_name=$(basename ${service_file} .service)
                systemctl disable --now ${service_name}
                rm -f ${service_file}
            fi
        done

        rm -rf ${APP_DIR}
        echo "uninstall ${APP_NAME} success"
    fi
}


echo "============================ ${APP_NAME} ============================"
echo "  1、install ${APP_NAME}"
echo "  2、update ${APP_NAME}"
echo "  3、uninstall ${APP_NAME}"
echo "======================================================================"
read -p "$(echo -e "place choose [1-3]：")" choose
case $choose in
1)
    install && wait && clear_install_files
    ;;
2)
    updateBinary && wait && clear_install_files
    ;;
3)
    uninstall && wait && clear_install_files
    ;;
*)
    echo "Input error, please try again!"
    ;;
esac