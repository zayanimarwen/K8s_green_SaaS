package auth

// Rôles disponibles dans la plateforme
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleViewer     = "viewer"
)

// RoleWeight permet de comparer les niveaux d'accès
var RoleWeight = map[string]int{
	RoleSuperAdmin: 100,
	RoleAdmin:      50,
	RoleViewer:     10,
}

// HasRole vérifie si une liste de rôles contient le rôle requis
func HasRole(userRoles []string, required string) bool {
	requiredWeight := RoleWeight[required]
	for _, r := range userRoles {
		if RoleWeight[r] >= requiredWeight {
			return true
		}
	}
	return false
}

// ExtractRoles extrait les rôles depuis les claims Keycloak
func ExtractRoles(realmAccess interface{}) []string {
	if realmAccess == nil {
		return nil
	}
	ra, ok := realmAccess.(map[string]interface{})
	if !ok {
		return nil
	}
	rolesRaw, ok := ra["roles"].([]interface{})
	if !ok {
		return nil
	}
	var roles []string
	for _, r := range rolesRaw {
		if s, ok := r.(string); ok {
			roles = append(roles, s)
		}
	}
	return roles
}
