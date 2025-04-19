package api

import "context"

func (a *ApiService) MailcowGetAlias(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/alias/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetDomain(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/domain/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetMailbox(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/mailbox/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetDkim(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/dkim/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetSyncjob(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/syncjobs/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetOAuth2Client(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/oauth2-client/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetOAuth2Clients(ctx context.Context) ApiMailcowGetAllRequest {
	return ApiMailcowGetAllRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/oauth2-client/all",
	}
}

func (a *ApiService) MailcowGetDomainAdmin(ctx context.Context, id string) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/alias/{id}",
		id:         id,
	}
}

func (a *ApiService) MailcowGetIdentityProviderKeycloak(ctx context.Context) ApiMailcowGetRequest {
	return ApiMailcowGetRequest{
		ApiService: a,
		ctx:        ctx,
		endpoint:   "/api/v1/get/identity-provider",
	}
}
