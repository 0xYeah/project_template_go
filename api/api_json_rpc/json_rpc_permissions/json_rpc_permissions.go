package json_rpc_permissions

var (
	APIPermissionTagDefault = "default"
)

type Permission struct {
	PermissionTag     string   `json:"permission_tag" comment:""`
	EncryptionEnabled bool     `json:"encryption_enabled"` // 是否开启加密
	AllowedUAs        []string `json:"allowed_uas"`
	AllowedMethods    []string `json:"allowed_methods"`
}
