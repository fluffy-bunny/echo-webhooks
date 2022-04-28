package auth

import (
	"echo-starter/internal/wellknown"

	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	middleware_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/middleware/oidc"
	services_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/claimsprincipal"
)

// BuildGrpcEntrypointPermissionsClaimsMap ...
func BuildGrpcEntrypointPermissionsClaimsMap() map[string]*middleware_oidc.EntryPointConfig {
	entryPointClaimsBuilder := services_claimsprincipal.NewEntryPointClaimsBuilder()
	// HEALTH SERVICE START
	//---------------------------------------------------------------------------------------------------
	// health check is open to anyone
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HealthzPath)

	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ReadyPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.HomePath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.AboutPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ErrorPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen("/static*")
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.ChannelPath)
	entryPointClaimsBuilder.WithGrpcEntrypointPermissionsClaimsMapOpen(wellknown.WebHookNoAuthPath)

	/*
			{
		  "aud": [
		    "b2b-client",
		    "users",
		    "invoices"
		  ],
		  "client_id": "b2b-client",
		  "exp": 1650817346,
		  "iat": 1650813746,
		  "iss": "https://echo-token-exchange.herokuapp.com/",
		  "jti": "c9immcnimfbc73quljt0",
		  "scope": [
		    "a",
		    "b",
		    "c",
		    "users.read",
		    "invoices"
		  ]
		}
	*/
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.WebHookPath).
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
		).GetChild().
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactTypeAndValue("aud", "invoices"),
		).GetChild().
		WithGrpcEntrypointPermissionsClaimFactsMapOR(
			services_claimsprincipal.NewClaimFactTypeAndValue("scope", "invoices"),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.WebHookBasicAuthPath).
		WithGrpcEntrypointPermissionsClaimFactsMapAND(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
			services_claimsprincipal.NewClaimFactTypeAndValue("auth_type", "basic"),
		)
	entryPointClaimsBuilder.GetClaimsConfig(wellknown.WebHookApiKeyPath).
		WithGrpcEntrypointPermissionsClaimFactsMapAND(
			services_claimsprincipal.NewClaimFactType(core_wellknown.ClaimTypeAuthenticated),
			services_claimsprincipal.NewClaimFactTypeAndValue("auth_type", "api-key"),
		)

	cMap := entryPointClaimsBuilder.GrpcEntrypointClaimsMap
	return cMap
}
