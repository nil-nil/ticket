package ticketjwt_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nil-nil/ticket/internal/domain"
	"github.com/nil-nil/ticket/internal/infrastructure/ticketjwt"
	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

var (
	publicKey  = []byte("-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAzGJoiXAk/6HGHL9YKlmf\npNGQByyCvHb0qtSgVu2BMmrWx5XWIDQNC2NHwucuvnwIkYq2G3aTsvrhtMdNtTCH\nh6f7TfORWv0lYBEkXfh16F4PA3pItbd4psUPc7EBtDpeGGfwOwnSH+dQ9c9zJPRB\nMawEa7HlmvlMkYtgzc6bh4MJfuKNBUX4b39pImJ2nDTbDvM2X9tLRnP80u8FFbEI\nAqsnVqLStUqFUKdMc9bOpitfJz0NVFEgeZ83ftD5xOvJuzu9Hz9P/xtdFDx1Rzim\n4or1kBvgHqaoa4CBO18DfW5nywv8wM8r8BuCCXDx+KBz34cESZ5nnlDmOHNOGARN\nox8PaQlSaGF0sIfrkH0X9///hSPZHFdMH876rvQnsEIawM6aqMg7usA9+UH0+NIb\nnhZ/8Z04WXSWojwGcP1jXjjYLJwF2T3B7tMRT8t4kC5in3RJMLs88TNaJUNu+HQH\nJQAxLZo1wWxWPldRxlfn1yabNF1Ite9AikjMcSXdB3Gg/N6Zb/2++omOhOWSn8tO\nBA9gYkKH7f0DlaB4Sdpu5CVKSeyebcqVe13HUBVWRoZ8JGYSDQNdPWn3P2ht2McJ\nvRvD7dNChRaA9+Eo7+wiX7w8tulCGzyAnBCuZaNXPxi3wMEdYPXmy2dhcU0rLO7b\n9VgbSpKx9vuZalKIWHKhxC0CAwEAAQ==\n-----END PUBLIC KEY-----\n")
	privateKey = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIJKAIBAAKCAgEAzGJoiXAk/6HGHL9YKlmfpNGQByyCvHb0qtSgVu2BMmrWx5XW\nIDQNC2NHwucuvnwIkYq2G3aTsvrhtMdNtTCHh6f7TfORWv0lYBEkXfh16F4PA3pI\ntbd4psUPc7EBtDpeGGfwOwnSH+dQ9c9zJPRBMawEa7HlmvlMkYtgzc6bh4MJfuKN\nBUX4b39pImJ2nDTbDvM2X9tLRnP80u8FFbEIAqsnVqLStUqFUKdMc9bOpitfJz0N\nVFEgeZ83ftD5xOvJuzu9Hz9P/xtdFDx1Rzim4or1kBvgHqaoa4CBO18DfW5nywv8\nwM8r8BuCCXDx+KBz34cESZ5nnlDmOHNOGARNox8PaQlSaGF0sIfrkH0X9///hSPZ\nHFdMH876rvQnsEIawM6aqMg7usA9+UH0+NIbnhZ/8Z04WXSWojwGcP1jXjjYLJwF\n2T3B7tMRT8t4kC5in3RJMLs88TNaJUNu+HQHJQAxLZo1wWxWPldRxlfn1yabNF1I\nte9AikjMcSXdB3Gg/N6Zb/2++omOhOWSn8tOBA9gYkKH7f0DlaB4Sdpu5CVKSeye\nbcqVe13HUBVWRoZ8JGYSDQNdPWn3P2ht2McJvRvD7dNChRaA9+Eo7+wiX7w8tulC\nGzyAnBCuZaNXPxi3wMEdYPXmy2dhcU0rLO7b9VgbSpKx9vuZalKIWHKhxC0CAwEA\nAQKCAgBe0kUXhlzT8HTeP3Oi6kOjjsoYWfEpiLLIUq80xSmMf363x+84r41hvCS9\n6s2H+DltdII6SZAKmFSAr3qA1kv6hteTea31Hb7qS+moYy4oqQKkJWZ4T+98x638\niaF1wSKIhigw68R9oq6v7BfKjDt21QyT/ku8025Pk+9MbE9B1mxgXrD4QlcZO25G\nUpIetHLxA6s1W6MXw5YHMncUcjZ6LneovQ+upi0llwhkMcNb2oFhzfRSKvU7F8AC\naOeIEbBc2kFKru/pNgO/8LCkY0chkUCOJDCdZ8p5XXwXVGRlHASxchVISpVi5xA0\nWx8XrzEzAgveL8x46aV1iSExCUYO02LC4lzAC2Ulr2yvS19ZYOUlJ/hB/3qPxKbK\nt0041vsuez1njt1cVpLiRrHUuI9lKdGx6UFdlqls6JAIeew8Hx0i2RM7n8IvKROQ\nmxfx2LPLBG9hSSg7XveZdCHgksBVsRqWBFmDuznZaVpxNvPURGv/WrBGZicgL5PU\nbMySWK3m+z6t8FAVETsXpOYfGxSmxppgxW8TTNzj+ieet1yFKBBHQMQLkhdF927w\nX9vMtqfXHcLdvDnGT0cvKYG4j4jbt0X/Au2/4WzcubRkar9eZPCqtMdMvMSwLRJ/\n1gAWUQTvg4vkW4+DYtRYHrBgw5S5S37WD8J+6pQZGMyd2mmCbQKCAQEA/d+/I6Qd\nb1yWDNbrFhQpV/omdZPPuYlHdF6Op011+H8drBdIyS2ZoYt6qai7+txcTTiuwnSH\n4CWkwuZ/02CCpCQZnSME11VrTJfV/shHmQOMh58d8YPhFBeQtr+m/mkgyxIskaY1\naAZHor5sq3UpuLn8KxkozbCYK1OCPE2auRYkQwxQCzNQ3B8vBar7vYV3C1Dnu6n/\ns27JvG1BdiuK9tXH+msYCKxoUDDRkyUs0v0AgT5AYipGUv80ziIo7Dda5ZQem5AE\nAIwFaXLHMgjFCgt3PMJsCnKD4uXT+h7YXabs6jP8eRSrIEVr8Vu/AzvS7COJk3MF\n0RiT85lh5m6G5wKCAQEAzhiQ9QQgum4cE2GC+VFp1lD6YFOLrolrwVmxIEmk2oKg\nfvQkLBwJdhD25g+aqGp6Oto6jZz4vU9w/bfbl+HAUKdLILcxbdfdfuQFW8JnNXL6\n+h2aHN8F1wt6eeGEzBrXuysV0nfjEenlJg6eLDdO+9MgewCNTed9mQw37krJTSoM\nTZrWwV24zhOc1e30nrRrriUWuL8O4WcwPpMjRptlmp5xfDrknOHOyfQWvORq6/8S\nEHmcshwijTXnyjz8Ue5rRSajVbt0uWbOFwK8JU84Mz+fw4SherRs6pjlyT/Xc8yu\nwVZLFw2peS75WhalxLbRbl3FM63KKl5RyYU7WDJ9ywKCAQBfL5D+z/2pT1GDJuGl\nuZF2xve8hdsQeYQtAXcDC2v76804RNKpe0tq4lzvV7CDcjO5UFNV3VNEm1iXKs0q\nd7kDyfVAkWyzP/enFBbMHFOb71S1VNdpQkUVv3Am1NzL7qa4/OtxAJxtkE0zm6oq\n7xbhh/ogPqKp3FhxhjICYiZs1vxplyg7ytW6coay3VTdzjjAKWl5V1fj8tn4qA4v\nPEwyGB2OqrCsL9g8mNE7FmkkAnA6BRkmtSsA22b1EqG9T1PpWAvRz8FwYw90ZfCB\ntgAKsBnY0hyoHh+M5xb/ZKlDE98oQK2cyD8RLnY2XGvVzoxatUhT3ICF0W1HnG60\ncyRpAoIBAG4uigjDrS+eQFpILnpWASw33LN01t93zmjvJ5foZz7+yQk2QsRmNNSv\nGyBBxWA2lKQ0GUuuWPj0qKasDbU0VtmHps2VwtJDrsHw68BzvTPBBdaDzumSfg/K\nri7M1287Boyk6yS7PWVNU1m2RO/EnfBZnirET8cPdIFHG/vEdbxQN4WhuyBjl7Js\nn9NrRPU35b2TTIN2eWEeBpfdl+VenMI95NQStDf/LMuhOCrCPztuAV1XduNt0TcH\nU6U2V3sB6M1ua7Ig5rVb9eAtcSLNKHGVmTcxCBeOsA/3sBmYjPn2upLYLIrlne4Q\no/R62SLCzlKfxRbs2YEvLbB8Dw8G52MCggEBAMzZCmJdoq00xAPAX7YRwY4tu8Eq\n1kNP7NcfROeOut1mBURZl0Zh6jJ0r5iEccpHN1bapyO4OuVfiMBz3IbaXr2wZ4Cg\nx0sDnGwWpvL28reY72/ke5chfIYT/o6tbXLXcnSqyRpnWo/9yeqy3IoYrU7bXyG+\n1mAnn5LgIyXTd8p9mmTjTLYQlj64FLS4PfZU9yzE4fYwRy4NnDVwJxvdFNQzQyk0\n+DMm6ayzkQy7pUkBqo7eh2kjGdBXGm3xaqIk8bp2gvYoH1GGvU8GI45OAlx8lXVF\n7HzyhGUQ+p9/drsSqGidDLNNH/kTd8w1Vp0PTEGKVlgiXguODfk+//7xZTM=\n-----END RSA PRIVATE KEY-----\n")
	userRepo   = mockUserRepository{
		users: map[uint64]domain.User{1: {ID: 1, FirstName: "firstname", LastName: "lastname"}},
	}
)

func TestIncorrectAlg(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")

	err = p.ValidateToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.")
	assert.ErrorIs(t, err, ticketjwt.ErrInvalidAlg)

	u, err := p.GetUser(context.Background(), "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.")
	assert.ErrorIs(t, err, ticketjwt.ErrInvalidAlg)
	assert.Equal(t, domain.User{}, u, "expect nil user when GetUser() errors")
}

func TestExpiredToken(t *testing.T) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": uint64(1),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(-1 * 24 * time.Hour).Unix(),
	})
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		assert.NoError(t, err, "mock key shouldn't error")
	}
	tokenString, err := token.SignedString(key)
	assert.NoError(t, err, "signing mock token shouldn't error")

	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")

	err = p.ValidateToken(tokenString)
	assert.ErrorIs(t, err, ticketjwt.ErrGettingToken)

	u, err := p.GetUser(context.Background(), tokenString)
	assert.ErrorIs(t, err, ticketjwt.ErrGettingToken)
	assert.Equal(t, domain.User{}, u, "expect nil user when GetUser() errors")
}

func TestUserSubjectFuncErr(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	_, err = p.NewToken(domain.User{})
	assert.EqualError(t, err, fmt.Sprintf("invalid jwt subject for user %+v", domain.User{}))
}

func TestInvalidPrivateKey(t *testing.T) {
	_, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, []byte{}, ticketjwt.RS512, 1000)
	assert.Error(t, err)
}

func TestInvalidPublicKey(t *testing.T) {
	_, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, []byte{}, privateKey, ticketjwt.RS512, 1000)
	assert.Error(t, err)
}

func TestNewTokenSuccess(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	_, err = p.NewToken(domain.User{ID: 1})
	assert.NoError(t, err, "successful userSubjectFunc should not error")
}

func TestValidateTokenSuccess(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	token, err := p.NewToken(domain.User{ID: 1})
	assert.NoError(t, err, "successful userSubjectFunc should not error")

	err = p.ValidateToken(token)
	assert.NoError(t, err)
}

func TestValidateTokenFailure(t *testing.T) {
	// Make the token expire immediately with 0 tokenLifetime
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 0)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	token, err := p.NewToken(domain.User{ID: 1})
	assert.NoError(t, err, "successful userSubjectFunc should not error")

	err = p.ValidateToken(token)
	assert.Error(t, err)
}

func TestValidateTokenError(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")

	err = p.ValidateToken("invalidtoken")
	assert.ErrorIs(t, err, ticketjwt.ErrGettingToken, "expected meaningful error")
}

func TestGetTokenUserSuccess(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserSuccessFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	token, err := p.NewToken(domain.User{ID: 1})
	assert.NoError(t, err, "valid user should not error")

	u, err := p.GetUser(context.Background(), token)
	assert.NoError(t, err)
	assert.Equal(t, userRepo.users[1], u)
}

func TestGetTokenUserFailure(t *testing.T) {
	p, err := ticketjwt.NewJwtAuthProvider(mockGetUserErrFunc, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	token, err := p.NewToken(domain.User{ID: 1})
	assert.NoError(t, err, "valid user should not error")

	u, err := p.GetUser(context.Background(), token)
	assert.Error(t, err)
	assert.Equal(t, domain.User{}, u)
}

func TestUserDeleted(t *testing.T) {
	userRepo.users[2] = domain.User{
		ID:        2,
		DeletedAt: ptr.To(time.Now()),
	}
	p, err := ticketjwt.NewJwtAuthProvider(userRepo.GetUser, publicKey, privateKey, ticketjwt.RS512, 1000)
	assert.NoError(t, err, "NewJwtAuthProvider should not error")
	token, err := p.NewToken(domain.User{ID: 2})
	assert.NoError(t, err, "valid user should not error")

	u, err := p.GetUser(context.Background(), token)
	assert.ErrorIs(t, err, ticketjwt.ErrUserDeleted)
	assert.Equal(t, domain.User{}, u)
}

type mockUserRepository struct {
	users map[uint64]domain.User
}

func (m *mockUserRepository) GetUser(ctx context.Context, userID uint64) (domain.User, error) {
	u, ok := m.users[userID]
	if !ok {
		return domain.User{}, fmt.Errorf("not found")
	}
	return u, nil
}

func mockGetUserSuccessFunc(ctx context.Context, userID uint64) (user domain.User, err error) {
	u, _ := userRepo.GetUser(ctx, 1)
	return u, nil
}

func mockGetUserErrFunc(ctx context.Context, userID uint64) (user domain.User, err error) {
	return domain.User{}, fmt.Errorf("mock error occurred")
}
