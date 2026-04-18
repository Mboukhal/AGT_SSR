package ep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Mboukhal/SvGoPg/cmd/settings"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	JWT_USER_KEY = "user"
)

func handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	// log.Println("Received magic link request for email:", email)
	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, `{"error": "Email not valid"}`, http.StatusBadRequest)
		return
	}

	if req.Password == "" || len(req.Password) < 6 {
		http.Error(w, `{"error": "Password must be at least 6 characters long"}`, http.StatusBadRequest)
		return
	}

	// hashedPassword, err := HashPassword(req.Password)
	// if err != nil {
	// 	http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
	// 	return
	// }

	// println("Registering user with email:", email, "and hashed password:", hashedPassword)

	// get user from db by email
	q := settings.GetQueries(r.Context())
	if q == nil {
		http.Error(w, `{"error": "Database queries not available"}`, http.StatusInternalServerError)
		return
	}

	user, err := q.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, `{"error": "User not found, please contact support!"}`, http.StatusNotFound)
		return
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	if user.Password != hashedPassword {
		http.Error(w, `{"error": "Incorrect password"}`, http.StatusUnauthorized)
		return
	}

	println("Found user with email:", user.ID.String(), "And Password:", user.Password)

	// generate JWT token with user ID and role as claims, expiring in 24 hours
	token, err := generateJWTToken(JWTClaims{
		UserID: user.ID.String(),
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})
	if err != nil {
		http.Error(w, `{"error": "Failed to generate JWT token"}`, http.StatusInternalServerError)
		return
	}

	// return success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("✓ User logged in with email %s.", email),
		"token":   token,
	})

}

func handleAuthRegister(w http.ResponseWriter, r *http.Request) {
	// For simplicity, using same handler for registration. In a real app, you'd have separate logic.
	var req AuthRequest

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, `{"error": "Email not valid"}`, http.StatusBadRequest)
		return
	}

	if req.Password == "" || len(req.Password) < 6 {
		http.Error(w, `{"error": "Password must be at least 6 characters long"}`, http.StatusBadRequest)
		return
	}

	// Here you would add logic to create the user in the database and hash the password

	err := CreateUser(r.Context(), email, req.Password, "")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to create user: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// hashedPassword, err := HashPassword(req.Password)
	// if err != nil {
	// 	http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
	// 	return
	// }

	// println("Registering user with email:", email, "and hashed password:", hashedPassword)

	// return success response as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("✓ User registered with email %s. You can now log in.", email),
	})
}

// // getUserRoleFromContext extracts the user role from JWT claims in context
// func GetUserRoleFromContext(ctx context.Context) string {
// 	claims, ok := ctx.Value(JWT_USER_KEY).(jwt.Claims)
// 	if !ok {
// 		return ""
// 	}
// 	// Try to extract role from custom claims
// 	// If you use map claims:
// 	if mapClaims, ok := claims.(jwt.MapClaims); ok {
// 		if role, ok := mapClaims["role"].(string); ok {
// 			return role
// 		}
// 	}
// 	// If you use a struct for claims, add logic here
// 	return ""
// }

func GetUserFromContext(ctx context.Context, key string) string {
	claims, ok := ctx.Value(JWT_USER_KEY).(jwt.Claims)
	if !ok {
		return ""
	}
	// Try to extract user ID from custom claims
	// log.Printf("Extracting user ID from context claims: %+v", claims)
	// If you use map claims:
	if mapClaims, ok := claims.(jwt.MapClaims); ok {
		if userID, ok := mapClaims[key].(string); ok {
			return userID
		}
	}
	// If you use a struct for claims, add logic here
	return ""
}

func AuthMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				//    h := w.Header()
				// h.Del("Content-Length")
				// h.Set("Content-Type", "text/plain; charset=utf-8")
				// h.Set("X-Content-Type-Options", "nosniff")
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintln(w, "UnAuthorized")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				// http.Error(w, "Invalid token format", http.StatusUnauthorized)
				// redirect to /login
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			// log.Printf("Validating token: %s", tokenString)

			// Validate token
			jwtSecret := os.Getenv("JWT_SECRET")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			//    log.Println("Authenticated user with token claims:", token.Claims)

			// If no roles specified, allow all authenticated users
			if len(allowedRoles) == 0 {
				// Token is valid, add claims to context
				ctx := context.WithValue(r.Context(), JWT_USER_KEY, token.Claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Check user role
			userRole := ""
			if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
				if role, ok := mapClaims["role"].(string); ok {
					userRole = role
				}
			}
			for _, role := range allowedRoles {
				if userRole == role {
					// Token is valid, add claims to context
					ctx := context.WithValue(r.Context(), JWT_USER_KEY, token.Claims)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

func HasRole(ctx context.Context, role string) bool {
	// claims, ok := ctx.Value(JWT_USER_KEY).(jwt.Claims)

	// println("Checking user role in HasRole function. Claims:", claims)

	// if !ok {
	// 	return false
	// }
	// // Try to extract role from custom claims
	// // If you use map claims:
	// if mapClaims, ok := claims.(jwt.MapClaims); ok {
	// 	println("Checking user role in HasRole function. Claims:", mapClaims)
	// 	if userRole, ok := mapClaims["role"].(string); ok {
	// 		println("Checking user role:", userRole, "against required role:", role)
	// 		return userRole == role
	// 	}
	// }
	ctx_role := GetUserFromContext(ctx, "role")
	if ctx_role == "" {
		println("User role not found in context")
		return false
	}
	if ctx_role == role {
		return true
	}
	println("User role in context:", ctx_role, "does not match required role:", role)
	return false
}

func handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	// return logout page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `
	<script>
		localStorage.removeItem('%s');
		window.location.href = '/';
	</script>
	`, os.Getenv("APP_TOKEN_ISSUER"))
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(JWT_USER_KEY).(jwt.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	q := settings.GetQueries(r.Context())
	if q == nil {
		http.Error(w, "Database queries not available", http.StatusInternalServerError)
		return
	}

	user, err := q.GetUserByEmail(r.Context(), GetUserFromContext(r.Context(), "email"))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	fmt.Printf("Authenticated user info: %+v\n", user)

	// Return user info as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"user_id":   user.ID.String(),
		"email":     user.Email,
		"username":  user.Username,
		"is_active": fmt.Sprintf("%t", user.IsActive),
	})
}

func RouterHandler(r chi.Router) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/ep/login", handleAuthLogin)
		r.Get("/logout", handleAuthLogout)
		r.Get("/me", handleMe)
		r.Route("/", func(r chi.Router) {
			r.Use(AuthMiddleware("ADMIN"))
			r.Post("/ep/register", handleAuthRegister) // for simplicity, using same handler for registration
			r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "Welcome to the protected admin route!")
			})
		})
	})
}
