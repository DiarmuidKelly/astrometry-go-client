# Branch Protection Rules

This document describes the branch protection rules configured for this repository.

## Main Branch Protection

The `main` branch is protected with the following rules:

### Required Reviews
- **Require pull request reviews before merging**: Enabled
- **Required approving reviews**: 1
- **Dismiss stale pull request approvals when new commits are pushed**: Enabled
- **Require review from Code Owners**: Enabled (see [CODEOWNERS](CODEOWNERS))

### Status Checks
- **Require status checks to pass before merging**: Enabled
- **Required status checks**:
  - `PR Title Validation` - Validates PR title format
  - `Go Tests` - Runs unit and integration tests
  - `Go Lint` - Runs golangci-lint

### Additional Restrictions
- **Require conversation resolution before merging**: Enabled
- **Require linear history**: Enabled (squash or rebase merges only)
- **Do not allow bypassing the above settings**: Enabled for administrators

### Allowed Merge Types
- ✅ **Squash merging** (recommended)
- ✅ **Rebase merging**
- ❌ **Merge commits** (disabled to maintain linear history)

### Auto-deletion
- **Automatically delete head branches**: Enabled

## Why These Rules?

### Pull Request Reviews
Ensures all code is reviewed by at least one code owner before merging, maintaining code quality and knowledge sharing.

### Status Checks
Automated checks catch issues before they reach the main branch:
- **PR Title Validation**: Ensures proper semantic versioning via PR titles
- **Tests**: Prevents broken code from merging
- **Linting**: Maintains code quality and consistency

### Linear History
Squash or rebase merging creates a clean, linear git history that's easier to understand and debug.

### Conversation Resolution
Ensures all review feedback is addressed before merging.

## Setting Up Branch Protection

To configure these rules on GitHub:

1. Go to **Settings** → **Branches**
2. Click **Add rule** under "Branch protection rules"
3. Set **Branch name pattern**: `main`
4. Configure the following options:

```
☑ Require a pull request before merging
  ☑ Require approvals (1)
  ☑ Dismiss stale pull request approvals when new commits are pushed
  ☑ Require review from Code Owners

☑ Require status checks to pass before merging
  ☑ Require branches to be up to date before merging
  Required status checks:
    - PR Title Validation
    - Go Tests
    - Go Lint

☑ Require conversation resolution before merging

☑ Require linear history

☑ Do not allow bypassing the above settings (includes administrators)
```

5. Click **Create** or **Save changes**

## Deploy Key Setup

The auto-release workflow requires a deploy key with write access:

1. Generate an SSH key pair:
   ```bash
   ssh-keygen -t ed25519 -C "github-actions-deploy-key" -f deploy_key -N ""
   ```

2. Add the public key (`deploy_key.pub`) to:
   **Settings** → **Deploy keys** → **Add deploy key**
   - Title: `github-actions-deploy-key`
   - Key: (paste contents of deploy_key.pub)
   - ☑ Allow write access

3. Add the private key (`deploy_key`) to:
   **Settings** → **Secrets and variables** → **Actions** → **New repository secret**
   - Name: `DEPLOY_KEY`
   - Value: (paste contents of deploy_key)

4. Delete the local key files:
   ```bash
   rm deploy_key deploy_key.pub
   ```

## Troubleshooting

### PR Can't Be Merged

**Issue**: "Required status checks have not been completed"
- **Solution**: Wait for all CI checks to complete. If they fail, fix the issues and push new commits.

**Issue**: "Review required"
- **Solution**: Request review from a code owner. See [CODEOWNERS](CODEOWNERS) for the list.

**Issue**: "Conversation not resolved"
- **Solution**: Resolve all comment threads before merging.

### Auto-Release Not Working

**Issue**: Release not created after PR merge
- **Solution**: Check that PR title follows the correct format (see [CONTRIBUTING.md](../CONTRIBUTING.md))

**Issue**: "Permission denied" during version bump commit
- **Solution**: Verify DEPLOY_KEY secret is set correctly with write permissions

## References

- [About branch protection rules](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [About code owners](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners)
- [About status checks](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/collaborating-on-repositories-with-code-quality-features/about-status-checks)
