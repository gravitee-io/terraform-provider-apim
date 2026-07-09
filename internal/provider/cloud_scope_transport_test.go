package provider

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCloudScopeTransportAllowsPermittedScope(t *testing.T) {
	claims := loadTestCloudClaims(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &http.Client{
		Transport: NewCloudScopeTransport(http.DefaultTransport, &claims),
	}

	req, err := http.NewRequest(http.MethodPut, server.URL+"/organizations/"+claims.Org+"/environments/"+claims.Envs[0]+"/apis", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCloudScopeTransportAllowsPermittedScopeOrgOnly(t *testing.T) {
	claims := loadTestCloudClaims(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &http.Client{
		Transport: NewCloudScopeTransport(http.DefaultTransport, &claims),
	}

	req, err := http.NewRequest(http.MethodPut, server.URL+"/organizations/"+claims.Org+"/settings", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestCloudScopeTransportRejectsUnauthorizedEnvironment(t *testing.T) {
	claims := loadTestCloudClaims(t)

	transport := NewCloudScopeTransport(http.DefaultTransport, &claims)
	req, err := http.NewRequest(http.MethodPut, "http://example.com/organizations/"+claims.Org+"/environments/unauthorized-env/apis", nil)
	require.NoError(t, err)

	_, err = transport.RoundTrip(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cloud token scope violation")
	assert.Contains(t, err.Error(), "unauthorized-env")
}

func TestCloudScopeTransportRejectsUnauthorizedOrgOnly(t *testing.T) {
	claims := loadTestCloudClaims(t)

	transport := NewCloudScopeTransport(http.DefaultTransport, &claims)
	req, err := http.NewRequest(http.MethodPut, "http://example.com/organizations/unauthorized-org/settings", nil)
	require.NoError(t, err)

	_, err = transport.RoundTrip(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cloud token scope violation")
	assert.Contains(t, err.Error(), "unauthorized-org")
}

func TestCloudScopeTransportRejectsUnauthorizedOrganization(t *testing.T) {
	claims := loadTestCloudClaims(t)

	transport := NewCloudScopeTransport(http.DefaultTransport, &claims)
	req, err := http.NewRequest(http.MethodPut, "http://example.com/organizations/unauthorized-org/environments/"+claims.Envs[0]+"/apis", nil)
	require.NoError(t, err)

	_, err = transport.RoundTrip(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cloud token scope violation")
	assert.Contains(t, err.Error(), "unauthorized-org")
}

func TestCloudScopeTransportPassthroughWithoutClaims(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	t.Cleanup(server.Close)

	client := &http.Client{
		Transport: NewCloudScopeTransport(http.DefaultTransport, nil),
	}

	req, err := http.NewRequest(http.MethodGet, server.URL+"/organizations/any/environments/any/apis", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusNoContent, res.StatusCode)
}

func TestCloudScopeTransportIgnoresPathsWithoutScope(t *testing.T) {
	claims := loadTestCloudClaims(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := &http.Client{
		Transport: NewCloudScopeTransport(http.DefaultTransport, &claims),
	}

	req, err := http.NewRequest(http.MethodGet, server.URL+"/health", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	require.NoError(t, res.Body.Close())
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestScopeFromRequestPath(t *testing.T) {
	orgEnv, ok := scopeFromRequestPath("/organizations/org-1/environments/env-1/apis/demo")
	assert.True(t, ok)
	assert.Equal(t, "org-1", orgEnv.Org)
	assert.Equal(t, "env-1", orgEnv.Env)

	_, ok = scopeFromRequestPath("/health")
	assert.False(t, ok)
}

func loadTestCloudClaims(t *testing.T) CloudTokenClaimsData {
	t.Helper()

	file, err := os.ReadFile("testdata/cloud-token.jwt")
	require.NoError(t, err)

	claims, err := extractCloudTokenData(string(file))
	require.NoError(t, err)
	return claims
}
