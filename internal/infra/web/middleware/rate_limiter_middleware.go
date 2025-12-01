package middleware

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/logger"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/configuration/rest_err"
	"github.com/Higor-ViniciusDev/posgo_raterlimite/internal/usecase/policy_usecase"
)

func RateLimiterMiddleware(policyUsecase policy_usecase.PolicyUsecaseInterface) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			input := policy_usecase.InputPolicyDTO{}

			apiKey := r.Header.Get("API-KEY")
			input.Tolken = apiKey

			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			input.IP = ip

			policy, key := policyUsecase.Resolver(r.Context(), input)

			logger.Info(fmt.Sprintf("Validando middleware ip: %v, tolken: %v", ip, apiKey))
			if err := policy.Validate(r.Context(), key); err != nil {
				restError := rest_err.ConvertInternalErrorToRestError(err)
				http.Error(w, restError.Message, restError.Code)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
