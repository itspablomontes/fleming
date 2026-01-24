package types

type PrincipalType string

const (
	PrincipalPatient    PrincipalType = "patient"
	PrincipalProvider   PrincipalType = "provider"
	PrincipalResearcher PrincipalType = "researcher"
	PrincipalSystem     PrincipalType = "system"
)

func ValidPrincipalTypes() []PrincipalType {
	return []PrincipalType{
		PrincipalPatient,
		PrincipalProvider,
		PrincipalResearcher,
		PrincipalSystem,
	}
}

func (pt PrincipalType) IsValid() bool {
	switch pt {
	case PrincipalPatient, PrincipalProvider, PrincipalResearcher, PrincipalSystem:
		return true
	}
	return false
}

type Principal struct {
	Address     WalletAddress   `json:"address"`
	Roles       []PrincipalType `json:"roles"`
	DisplayName string          `json:"displayName,omitempty"`
}

func NewPrincipal(address WalletAddress, roles ...PrincipalType) (Principal, error) {
	if address.IsEmpty() {
		return Principal{}, ErrInvalidAddress
	}

	if len(roles) == 0 {
		return Principal{}, NewValidationError("roles", "at least one role is required")
	}

	for _, role := range roles {
		if !role.IsValid() {
			return Principal{}, NewValidationError("roles", "invalid principal role: "+string(role))
		}
	}

	return Principal{
		Address: address,
		Roles:   roles,
	}, nil
}

func (p Principal) HasRole(t PrincipalType) bool {
	for _, role := range p.Roles {
		if role == t {
			return true
		}
	}
	return false
}

func (p Principal) IsPatient() bool {
	return p.HasRole(PrincipalPatient)
}

func (p Principal) IsProvider() bool {
	return p.HasRole(PrincipalProvider)
}

func (p Principal) CanOwn() bool {
	return p.IsPatient()
}

func (p Principal) CanGenerate() bool {
	return p.IsProvider() || p.HasRole(PrincipalResearcher)
}
