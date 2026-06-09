# Gitea 1.23 Feature Candidates

> Criteria: (1) Appears in v1.23.0 release notes FEATURES section, (2) Has a linked GitHub issue, (3) PR was properly merged into v1.23.0.

## Summary

- **Total FEATURES in release notes:** 25
- **Meeting all 3 criteria:** 19
- **Missing linked issue:** 9 (2 are base-module PRs whose companion PRs do have issues)

---

## Features Meeting All Criteria (19)

| # | Feature | PR | Issue | Issue Title | Files Changed | Lines Changed |
|---|---------|-----|-------|-------------|---------------|---------------|
| 1 | Add Arch package registry | [#32692](https://github.com/go-gitea/gitea/pull/32692) | [#25037](https://github.com/go-gitea/gitea/issues/25037) | Arch linux packages | 43 | +1687 -91 |
| 2 | Allow to disable the password-based login (sign-in) form | [#32687](https://github.com/go-gitea/gitea/pull/32687) | [#7633](https://github.com/go-gitea/gitea/issues/7633) | Is it possible to choose default auth source? | 7 | +73 -48 |
| 3 | Allow to fork repository into the same owner | [#32819](https://github.com/go-gitea/gitea/pull/32819) | [#22882](https://github.com/go-gitea/gitea/issues/22882) | Fork your own repository | 5 | +44 -5 |
| 4 | Allow cropping an avatar before setting it | [#32565](https://github.com/go-gitea/gitea/pull/32565) | [#31990](https://github.com/go-gitea/gitea/issues/31990) | Unable to display complete avatar | 12 | +80 -9 |
| 5 | Support "merge upstream branch" (Sync fork) | [#32741](https://github.com/go-gitea/gitea/pull/32741) | [#20880](https://github.com/go-gitea/gitea/issues/20880) | Single-action web UI for a fork to fast-forward pull from origin (Sync fork) | 10 | +323 -136 |
| 6 | Add reviewers' selection to new pull request | [#32403](https://github.com/go-gitea/gitea/pull/32403) | [#26289](https://github.com/go-gitea/gitea/issues/26289) | Assign reviewers on pull request creation | 26 | +500 -268 |
| 7 | Suggestions for issues | [#32327](https://github.com/go-gitea/gitea/pull/32327) | [#16872](https://github.com/go-gitea/gitea/issues/16872) | autocompletion / suggestion for # (reference issues/pulls) | 9 | +202 -48 |
| 8 | Included tag search capabilities | [#32045](https://github.com/go-gitea/gitea/pull/32045) | [#31998](https://github.com/go-gitea/gitea/issues/31998) | Add searching capabilities to the tags page | 4 | +33 -7 |
| 9 | Add option to filter board cards by labels and assignees | [#31999](https://github.com/go-gitea/gitea/pull/31999) | [#21846](https://github.com/go-gitea/gitea/issues/21846) | Project board card filtering | 14 | +325 -33 |
| 10 | Introduce globallock as distributed locks | [#31813](https://github.com/go-gitea/gitea/pull/31813) | [#19620](https://github.com/go-gitea/gitea/issues/19620) | replace sync module | 13 | +185 -107 |
| 11 | Support compression for Actions logs & enable by default | [#32013](https://github.com/go-gitea/gitea/pull/32013) | [#31801](https://github.com/go-gitea/gitea/issues/31801) | Enable compression for Actions logs by default | 2 | +3 -3 |
| 12 | Add pure SSH LFS support | [#31516](https://github.com/go-gitea/gitea/pull/31516) | [#17554](https://github.com/go-gitea/gitea/issues/17554) | Support LFS purely over SSH protocol | 13 | +945 -53 |
| 13 | Add Passkey login support | [#31504](https://github.com/go-gitea/gitea/pull/31504) | [#22015](https://github.com/go-gitea/gitea/issues/22015) | Add support for passkeys (WebAuthn as primary authentication) | 8 | +184 -11 |
| 14 | Actions support workflow dispatch event | [#28163](https://github.com/go-gitea/gitea/pull/28163) | [#23668](https://github.com/go-gitea/gitea/issues/23668) | Actions - Manually trigger a workflow/action | 10 | +580 -17 |
| 15 | Support repo license | [#24872](https://github.com/go-gitea/gitea/pull/24872) | [#278](https://github.com/go-gitea/gitea/issues/278) | Display a License tab | 47 | +906 -22 |
| 16 | Issue time estimate, meaningful time tracking | [#23113](https://github.com/go-gitea/gitea/pull/23113) | [#23112](https://github.com/go-gitea/gitea/issues/23112) | Issue time estimate, meaningful time tracking | 21 | +390 -164 |
| 17 | Rearrange Clone Panel | [#31142](https://github.com/go-gitea/gitea/pull/31142) | [#23202](https://github.com/go-gitea/gitea/issues/23202) | Clone Button Rearrangement | 19 | +191 -195 |
| 18 | Enhancing Gitea OAuth2 Provider with Granular Scopes | [#32573](https://github.com/go-gitea/gitea/pull/32573) | [#31609](https://github.com/go-gitea/gitea/issues/31609) | Enhancing Gitea OAuth2 Provider with Granular Scopes for Resource Access | 8 | +537 -18 |
| 19 | Use env GITEA_RUNNER_REGISTRATION_TOKEN as global runner token | [#32946](https://github.com/go-gitea/gitea/pull/32946) | [#23703](https://github.com/go-gitea/gitea/issues/23703) | Improve Config Management/Stateless Runner Deploy Workflows | 8 | +152 -18 |

## Features NOT Meeting All Criteria (9)

All PRs were merged into v1.23.0, but no linked GitHub issue was found.

| # | Feature | PR | Note |
|---|---------|-----|------|
| 1 | Support quote selected comments to reply | [#32431](https://github.com/go-gitea/gitea/pull/32431) | No linked issue |
| 2 | Add priority to the protected branch | [#32286](https://github.com/go-gitea/gitea/pull/32286) | No linked issue |
| 3 | Add automatic light/dark option for the colorblind theme | [#31997](https://github.com/go-gitea/gitea/pull/31997) | No linked issue |
| 4 | Support migration from AWS CodeCommit | [#31981](https://github.com/go-gitea/gitea/pull/31981) | No linked issue |
| 5 | GitHub like repo home page | [#32213](https://github.com/go-gitea/gitea/pull/32213) | #27931 is a PR, not an issue |
| 6 | Tweak repo sidebar | [#32847](https://github.com/go-gitea/gitea/pull/32847) | No linked issue, companion to #32213 |
| 7 | Update i18n.go - Language Picker | [#32933](https://github.com/go-gitea/gitea/pull/32933) / [#32935](https://github.com/go-gitea/gitea/pull/32935) | No linked issue |
| 8 | Introduce globallock (base module) | [#31908](https://github.com/go-gitea/gitea/pull/31908) | Companion to #31813 which has issue #19620 |
| 9 | Support compression for Actions logs (base module) | [#31761](https://github.com/go-gitea/gitea/pull/31761) | Companion to #32013 which has issue #31801 |
