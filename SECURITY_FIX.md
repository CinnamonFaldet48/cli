# Prevent silent fallback to unauthenticated requests

**Issue:** When TokenForHost returns error (e.g., keyring locked), ApiRun previously fell back to unauthenticated request, leaking unauthenticated calls and hiding the error.

**Fix:** ResolveToken now surfaces the error and ApiRun returns it immediately with host context: "failed to get authentication token for host X: keyring ...". Empty token now returns "no authentication token found..." instead of silent fallback. Ensures strict token verification.

**Tests:** TestApiRun_TokenResolutionErrorSurfaced and TestApiRun_MissingTokenReturnsError both pass.
