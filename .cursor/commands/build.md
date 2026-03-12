Build all workspaces and verify the build succeeds.

1. Run bunx turbo build to build all workspaces
2. If build fails, diagnose the error from the output
3. Fix any build issues (delegate to @hephaestus if complex)
4. Re-run until build passes with exit code 0
5. Report results: which workspaces built successfully, any warnings

For individual workspaces:
- Frontend only: bunx turbo build --filter=@axiom/web
- Simulation only: bunx turbo build --filter=@axiom/sim
