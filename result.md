# Goal
Backport branches and tags that were only maintained on `github.com/docker/docker-ce` to the 3 components `github.com/docker/cli`, `github.com/docker/engine` and `github.com/docker/docker-ce-packaging`.
The following branches have to be backported:
* 17.06
* 17.07
* 17.09
* 17.10
* 17.11
* 17.12
* 18.01
* 18.02
* 18.03
* 18.04
* 18.05

Starting **18.06** the branches were also maintained on each component (`18.06`, `18.09` and `19.03`).

On docker-ce a bot was merging each pull request on each component in the same repository under the `components` directory. But some commits were made directly on each branches and never backported to the components. We need to identify them and cherrypick the commits from docker-ce to the component repositories.

Tags were done on docker-ce and components, but only a subset of tags were backported to components. We need to identify them and find the right commit on components to tag.

# Process

## Branches

1/ First step is using `git filter-branch` on docker-ce repository for each components (using `components/[component] directoy) for each missing branches to extract the .
```
$ git filter-branch -f --prune-empty --subdirectory-filter components/engine  HEAD -- --all
```

The result is a translation branch in docker-ce with only the commits targeting the `components/[component]` directory and the root directory became the component directory.

Example:
```
components/engine/VERSION -> VERSION
```

:warning: The main **drawback** here is that every commit is rewritten so the hash will differ from the original commit.

2/ On docker-ce repository we now need to find the first common ancestor between a specific untranslated branch (eg `18.05`) and the master branch. Once we find it, we will rewind history from first parent to first parent until we find a commit targeting the component.

This commit comes from the merge bot and has the form "Merge component 'component'..."

We fetch its last commit parent (the merge commit of the pull request) and get its message.


```
On docker-ce repository:

master  18.05
  |       |
  |      /
  |    /
  |  /
  common ancestor
  |
  some commits
  |
  merge bot commit "Merge component 'engine'..."
  |   \
  |    \
  |     "Merge pull request #XX ...."
```

3/ We now start from the translated/extracted component branch and rewind until we find the commit comparing the messages.
```
18.05-extract-engine
  |
  |
  |
  |
  |
  "Merge pull request #XX"
```
We have a beginning (head of the branch) and an end (the merge pull request).
We can crawl back all the commits to cherry pick starting from the head and store the hashes somewhere.
```
5cfd81df09fa69be905d8721e17d427c7ced709c
1163999891fb6fbc1648d4288d72007a99d890c4
10a4ab117cbfca142455a21fb49dc235175b2ca2
604aafeefd12be2ba97845d49b5289d757e3b119
c56f3d934e80fb5bbb26cc433753fdd87d81a521
...
```

4/ Now we move to the component repository (eg: `docker/engine`). Starting from master, we rewind the history until we find the merge PR commit (comparing the commit messages) and branch from here.
```
On docker/engine:

master
  |
  |
  |
  |
  "Merge pull request #XX" <- we branch 18.05 from here
```

5/ We cherry pick all the commits we stored, with reversed order.

```
On docker/engine:

master    18.05
  |      / cherry pick 5cfd81df09fa69be905d8721e17d427c7ced709c
  |     / ...
  |    / cherry pick 10a4ab117cbfca142455a21fb49dc235175b2ca2
  |   / cherry pick 604aafeefd12be2ba97845d49b5289d757e3b119
  |  / cherry pick c56f3d934e80fb5bbb26cc433753fdd87d81a521
  | /
  "Merge pull request #XX"
```

:warning: As everything was flattened, commit order may not be "guaranteed", but if we diff the last commit of the branches, diff is empty.
It's also guaranteed that there wasn't any conflict during the process.

6/ Last step, we verify that the end result branch has no diff with the translated branch from docker-ce.

```
git diff 18.05-extract-engine 18.05
# nothing
```

# Results

## Branches
### 17.06
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.06-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.06-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.06-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.06-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.06-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.06-extract-packaging)**

### 17.07
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.07-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.07-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.07-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.07-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.07-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.07-extract-packaging)**

### 17.09
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.09-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.09-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.09-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.09-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.09-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.09-extract-packaging)**

### 17.10
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.10-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.10-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.10-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.10-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.10-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.10-extract-packaging)**

### 17.11
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.11-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.11-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.11-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.11-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.11-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.11-extract-packaging)**

### 17.12
CLI: https://github.com/silvin-lubecki/cli-extract/commits/17.12-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.12-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/17.12-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.12-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/17.12-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/17.12-extract-packaging)**

### 18.01
CLI: https://github.com/silvin-lubecki/cli-extract/commits/18.01-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.01-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/18.01-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.01-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/18.01-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.01-extract-packaging)**

### 18.02
CLI: https://github.com/silvin-lubecki/cli-extract/commits/18.02-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.02-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/18.02-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.02-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/18.02-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.02-extract-packaging)**

### 18.03
CLI: https://github.com/silvin-lubecki/cli-extract/commits/18.03-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.03-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/18.03-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.03-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/18.03-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.03-extract-packaging)**

### 18.04
CLI: https://github.com/silvin-lubecki/cli-extract/commits/18.04-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.04-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/18.04-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.04-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/18.04-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.04-extract-packaging)**

### 18.05
CLI: https://github.com/silvin-lubecki/cli-extract/commits/18.05-extract-cli **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.05-extract-cli)**

Engine: https://github.com/silvin-lubecki/engine-extract/commits/18.05-extract-engine **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.05-extract-engine)**

Packaging: https://github.com/silvin-lubecki/packaging-extract/commits/18.05-extract-packaging **[[compare]](https://github.com/silvin-lubecki/docker-ce/commits/18.05-extract-packaging)**

### For the record
You will find [here](https://github.com/silvin-lubecki/docker-ce) a fork of docker/docker-ce with branches for each components and each version, with all the extracts using the following command
```
$ git filter-branch -f -d /tmp/rebase --prune-empty --subdirectory-filter components/engine  HEAD -- --all
```
As this command is computationally heavy (~45mn for each branch on `docker/engine`), this repo can be used as an intermediary step if we need to compute branches again.

Ex:
Extract of docker-ce 18.05 for docker engine: https://github.com/silvin-lubecki/docker-ce/tree/18.05-extract-engine

Extract of docker-ce 17.12 for docker cli: https://github.com/silvin-lubecki/docker-ce/tree/17.12-extract-cli

Extract of docker-ce 17.06 for docker packaging: https://github.com/silvin-lubecki/docker-ce/tree/17.06-extract-packaging

## Tags

We also have some tags that were never backported to each components. We need to identify those tags, and find the right commit to tag on each component.

### Engine

| Tag | Docker-ce sha | docker/engine sha |
| --- | --- | --- |
| v17.06.0-ce | 02c1d876176546b5f069dae758d6a7d2ead6bd48 | bda340abf16b691feac794480cbcce2c711f31fc |
| v17.06.0-ce-rc1 | 2bcfe6ffc2eda302f8a137e996755ec49b57d17d | 3240dfbec53923d05f4b07a9212178a63c3fdd3a |
| v17.06.0-ce-rc2 | 402dd4a9ea3802e45718f871b00efe52caae5108 | 85a3b3402b0c621474a37856c48a0140ac4e3115 |
| v17.06.0-ce-rc3 | 7953dbc64969ce8e0e13d77a0793ead46456feaf | e0969deefba15c4ca436b27b569b16e1f96b31b2 |
| v17.06.0-ce-rc4 | 0f80f28c1a66ce092f5cd75ddbb566d1301d3dfa | d9437def7d61380e32514e00bf0b5c89b1745a7e |
| v17.06.0-ce-rc5 | b7e417335c2f5b1b8ec73ac5f6a340a0ecc7367f | 84c2adb399a82eb1d9bb6bc1fe98dc499bf38476 |
| v17.06.1-ce | 874a7374f31c77aca693d025101b2de1b20b96c2 | 81f63391e481db17028122f3bbca6271a57c0e8d |
| v17.06.1-ce-rc1 | c257617b99f931ff36128557d186f1041f6b8d44 | 7b33830d6e26e863520dc6fba5103f44b6316ec4 |
| v17.06.1-ce-rc2 | 96d84243ab0089d81353b5c0931a5e512699bd58 | 5d71d08d0f2f5c5261c2bf76127882b848710cad |
| v17.06.1-ce-rc3 | 344dbf919d8e52fbadbbf984d8c640bd375cf7c4 | b3f242d00082b3dae096eb7aecfb9fe566ec5614 |
| v17.06.1-ce-rc4 | 0a33916f54b9014f52a74e4c7d5bee7084acfe0f | 8e1f2c3fd970aec9a7243d45cf51fcf444a9ffc3 |
| v17.06.2-ce | cec0b72a9940e047e945a09e1febd781e88366d6 | eed47e15052b7edba92122459faf9c28cabeb2c1 |
| v17.06.2-ce-rc1 | 34d73cd0d87de014897165c339a5c4f1f13718df | e36a7624e748b6b5b80b8227d33992de9f46e79e |
| v17.07.0-ce | 87847530f7176a48348d196f7c23bbd058052af1 | 6e7a9e1d30ca3715a22842f9320cf82662d43203 |
| v17.07.0-ce-rc1 | 8c4be39ddda5ae43c54ca89780a36c6ecfda5117 | 5a521cfa48a54a68443cb1c0790e2b0a682628bd |
| v17.07.0-ce-rc2 | 36ce6055c2545dcb386e054ac835f28b91f33364 | 562b165bc9e43350253422a4d960fcf24b9f19b4 |
| v17.07.0-ce-rc3 | 665d244001888bcc73c6b2ae5ba486ca7b61340c | 1062b84197aa7bd6b8e99ecc7aee5d5b5e9e8929 |
| v17.07.0-ce-rc4 | fe143e3bbc09eb32c33dfe1a86cc41bf11d4d2a9 | 73dd43eb19bdd6102867f505f7a40704727acae8 |
| v17.09.0-ce | afdb6d44a80f777069885a9ee0e0f86cf841b1bb | 25824859cf8bfebd2ec6f94c619ec40593daf5af |
| v17.09.0-ce-rc1 | ae21824df8b1bb712da98513788231dcc0643ccc | f903137ce1ecab819afe4df0eb258f53793ef874 |
| v17.09.0-ce-rc2 | 363a3e75cdf1d11ad3a41ec2585dba1009f55e76 | 13bfc729718ea6ee1232f0c8bfc0ef7ed6cb55ba |
| v17.09.0-ce-rc3 | 2357fb28b503963f2313d407bc4af52476512905 | 86925fd03ed36ae21470e7c2d3596b5fdc93bee9 |
| v17.09.1-ce | 19e2cf6259bd7f027a3fff180876a22945ce4ba8 | a2f00c813c52f45fa856d3711e8023a981630b40 |
| v17.09.1-ce-rc1 | 2d63290f4bc928dbdf5425d646b27a269f813cd1 | 987ccc166c44a5710073ea022553de824f0ce4a2 |
| v17.10.0-ce | f4ffd2511ce93aa9e5eefdf0e912f77543080b0b | de0545be0b5d4c3870d98c3df01a94ac484cd947 |
| v17.10.0-ce-rc1 | d866876d86e5228e367fca3318c7c0fd3fbb1ea9 | 8db75cb6f4bd361a1526fb10c4310b7d79648ca9 |
| v17.10.0-ce-rc2 | af94197a3641aeaa3c3aaf9f6ca994b87e7382c6 | 65fd4bcd513008cd6f735687eafbaa008c98d6ec |
| v17.11.0-ce | 1caf76ce6baa889133ece59fab3c36aaf143d4ef | dfe625f3b8a8bdb2243d386fbe8a15bb34364005 |
| v17.11.0-ce-rc1 | df7b627909ca07c0a1a000fde37b9a7f264f0e6a | fb9040837db666308ee4b111168fc5f119bfc09f |
| v17.11.0-ce-rc2 | d7062e5443166bb01b04d999347d22d742fdf848 | 356f701725de652f213e57cd58129cdbbd8dfa46 |
| v17.11.0-ce-rc3 | 5b4af4f7126e185917784eeac63b89b1a03bc9bb | 0d5c4e2727aeee3081f23a64bbbbf45cc2734747 |
| v17.11.0-ce-rc4 | 587f1f003bbe3090710151cc9e89042959d55d0b | 253d452fc30ac3221aa206ed89e1be5433eb2f6b |
| v17.12.0-ce | 740d71bf6a4ae71c667645903a5fac08568ca830 | 8315e9587e12dacabaa5e6dccdb5426f1166485d |
| v17.12.0-ce-rc1 | 704b884e52487c84167372654d313aacb64338f2 | a023a599913439f0a08adffc3f242ce187fd8bdd |
| v17.12.0-ce-rc2 | 184c27d7854e74a1a1541d8aa6e3c04a641550a0 | dd807791ff13616681943c4d9cf171e5f528cb8c |
| v17.12.0-ce-rc3 | 3083961ff088e5bc70512657942391f6fe7f8daa | dcf817bf2b9f78f0b264398285dae90587e34a14 |
| v17.12.0-ce-rc4 | 740d71bf6a4ae71c667645903a5fac08568ca830 | 8315e9587e12dacabaa5e6dccdb5426f1166485d |
| v17.12.1-ce | 319872e09a8f6c69844f7637f9398ae27b230dde | 08bfa1cff3927e4209b8a1e738138b5709534eaa |
| v17.12.1-ce-rc1 | 602216ce5634bf938d194d0546c5b584de2925cb | 21968c29642a95a4508def6f8b3729384332ec1d |
| v17.12.1-ce-rc2 | 319872e09a8f6c69844f7637f9398ae27b230dde | 08bfa1cff3927e4209b8a1e738138b5709534eaa |
| v18.01.0-ce | 4dd1acceb629c776b48568b3a80760d7fd1917df | 7e3cce37eda990f9eec184890647dd40d9a1c074 |
| v18.01.0-ce-rc1 | 4dd1acceb629c776b48568b3a80760d7fd1917df | 7e3cce37eda990f9eec184890647dd40d9a1c074 |
| v18.02.0-ce | f25c14edf2b30682417c65a7f5ebc987cd4cef50 | 84e6b4c6799262cfb860da0c9eb767c016ea23f5 |
| v18.02.0-ce-rc1 | 3b28f86d447e71b0ce222c69120e7ea088bff527 | f20708bda0a803d1b4cf3b1a34f13617109cbb47 |
| v18.02.0-ce-rc2 | f25c14edf2b30682417c65a7f5ebc987cd4cef50 | 84e6b4c6799262cfb860da0c9eb767c016ea23f5 |
| v18.03.0-ce | 78455c2b2ff210bb7ff94d41b13e2cf53b76c1e0 | feb479ed66340f7016f4e5de5ec2741dafe97e91 |
| v18.03.0-ce-rc1 | 5c06a61da41ba95a370e878ea760b0e9e4c09d51 | c85afd025323d817254f3d678128b6ec37ae7158 |
| v18.03.0-ce-rc2 | 5ba2b1a74d741c97f5627f281416a09fd517fa06 | aba2685ee0104c71c65037f644df3406cf4daabf |
| v18.03.0-ce-rc3 | 23a90170376e3319697c8d5477c91f3308ce6299 | 09ddb49879eaa1dab68aef6a4f5f508732d8fa72 |
| v18.03.0-ce-rc4 | 78455c2b2ff210bb7ff94d41b13e2cf53b76c1e0 | feb479ed66340f7016f4e5de5ec2741dafe97e91 |
| v18.03.1-ce | 6d0d01b238330adf8179d9485316e8e4c6ec46cd | d38a8fd8dc8bb9a09ae3b6304104f531cfbc1c1a |
| v18.03.1-ce-rc1 | f10cbe710b9476f20135c03a4c43879c47a880ad | c5cdf5c787fd229e5df079ccce11c898081568ca |
| v18.03.1-ce-rc2 | af0331718036f483dc51dd154cec89fb1a850b36 | 2ff4e313c7ff61c60c7c4f6eee4b1903b7111a59 |
| v18.04.0-ce | 788bd7e0de8ff0ef76faec170b08c78f317d7bbc | 8f305b11fe8073f98bac80bf796bc97f8501ca9b |
| v18.04.0-ce-rc1 | 8a9fee12aedbc2eed379c04478c74062c551c33f | ed7b6428c133e7c59404251a09b7d6b02fa83cc2 |
| v18.04.0-ce-rc2 | 788bd7e0de8ff0ef76faec170b08c78f317d7bbc | 8f305b11fe8073f98bac80bf796bc97f8501ca9b |
| v18.05.0-ce | 1b9a860800a134de2fc6f6e27ed8660033bb99c5 | 05fb243667d76c6861130979f719ee2d1a810ab9 |
| v18.05.0-ce-rc1 | 1b9a860800a134de2fc6f6e27ed8660033bb99c5 | 05fb243667d76c6861130979f719ee2d1a810ab9 |
| v18.09.9-beta1 | 0950f0f1fa79be1405d31ba24bbbe992bc4944a5 | abda7e1227ba45584a6abbf6fbc07e31f970380f |
### CLI

| Tag | Docker-ce sha | docker/cli sha |
| --- | --- | --- |
| v17.06.0-ce | 02c1d876176546b5f069dae758d6a7d2ead6bd48 | 57d3a0e1dad54d77594cbba23a6418ecbf5ac230 |
| v17.06.0-ce-rc1 | f3810787c8c7d4aa21820f78e17c86ec701bf7cc | b4864fad111a663abf8c9d6991857d9fd060cbbf |
| v17.06.0-ce-rc2 | 402dd4a9ea3802e45718f871b00efe52caae5108 | 8c7aba4ab7725aed13a52b77dd8d33c88970e0d6 |
| v17.06.0-ce-rc3 | 7953dbc64969ce8e0e13d77a0793ead46456feaf | 0ac54d58f8b35e0ba5c0ebd91f5829716ac1d6f6 |
| v17.06.0-ce-rc4 | 0f80f28c1a66ce092f5cd75ddbb566d1301d3dfa | dcd5c998c6c2d96ff5ed05312333382d3e7820bb |
| v17.06.0-ce-rc5 | b7e417335c2f5b1b8ec73ac5f6a340a0ecc7367f | 7ba7f35d423e43bc9cb5b5a387d1a935a579b88c |
| v17.06.1-ce | 874a7374f31c77aca693d025101b2de1b20b96c2 | ea0a0373c441da65d40e5816d287d9d0c6d47581 |
| v17.06.1-ce-rc1 | c257617b99f931ff36128557d186f1041f6b8d44 | 3c336dc761730e3dd95fdbbc258d4687f66d8a43 |
| v17.06.1-ce-rc2 | 96d84243ab0089d81353b5c0931a5e512699bd58 | 939313f4ea8851e81a423c5240f3682620373492 |
| v17.06.1-ce-rc3 | 87a7a648dfa3e41990fc097c01ee1a814b5fdd0d | 5f9e8c298d374e7fce2db5547555cfd7fb51872e |
| v17.06.1-ce-rc4 | 0a33916f54b9014f52a74e4c7d5bee7084acfe0f | 311e8aa535e8df2a99a53b1c414038387f57924d |
| v17.06.2-ce | cec0b72a9940e047e945a09e1febd781e88366d6 | 03b34f1f5faa65c3e60ab1d7a0d2ecf99e804c9e |
| v17.06.2-ce-rc1 | 34d73cd0d87de014897165c339a5c4f1f13718df | 1e9dba4819dc14d10c3c6e24e3567551d1494237 |
| v17.07.0-ce | 87847530f7176a48348d196f7c23bbd058052af1 | d74686e6c1a49497998b38ff0f344250683c32a8 |
| v17.07.0-ce-rc1 | 8c4be39ddda5ae43c54ca89780a36c6ecfda5117 | 161a76157469569ff87ca07b498e85a3a6ca1969 |
| v17.07.0-ce-rc2 | 82a59cdb2c53db58212296be3de0e6eaa3724b6b | 807254cf0b9ce863368dcf2e772f585fcd217de3 |
| v17.07.0-ce-rc3 | 3f7ca09fc65cd8d9b9e3b37229d555fe44c62292 | 60c9cf8eb8d1bbf840f4c5db96faef32b2789c1e |
| v17.07.0-ce-rc4 | fe143e3bbc09eb32c33dfe1a86cc41bf11d4d2a9 | 381f6dbc022389e1ebea7c4c9e442dda765288ed |
| v17.09.0-ce | afdb6d44a80f777069885a9ee0e0f86cf841b1bb | 551c0b2380c26805c47aa4a3339109d5754d3493 |
| v17.09.0-ce-rc1 | ae21824df8b1bb712da98513788231dcc0643ccc | 5a339b6b32747f14e2c94c5aa2788bd1bb75c520 |
| v17.09.0-ce-rc2 | 17db25d09805ba7f4d0d42693f5738b826d7c607 | 48e2639411c651942b01c8c4e19676a8b512accc |
| v17.09.0-ce-rc3 | 2357fb28b503963f2313d407bc4af52476512905 | eea938447625c5a09fcf56ef04ec99535563da8b |
| v17.09.1-ce | 19e2cf6259bd7f027a3fff180876a22945ce4ba8 | 3fb95133411fac45cc5b278e723264efcf52946e |
| v17.09.1-ce-rc1 | 2d63290f4bc928dbdf5425d646b27a269f813cd1 | 5fac918b829fed534ad66e14a2e1e1de3c125e0e |
| v17.10.0-ce | f4ffd2511ce93aa9e5eefdf0e912f77543080b0b | b524db30ced63fbb3ec1362d031cfb59b06800a1 |
| v17.10.0-ce-rc1 | d866876d86e5228e367fca3318c7c0fd3fbb1ea9 | d2cf71ce5c005935d5833a5b1f78a7e5a8c31530 |
| v17.10.0-ce-rc2 | af94197a3641aeaa3c3aaf9f6ca994b87e7382c6 | fc8d6a29197c9535cb8c66c32b9414732fd28b05 |
| v17.11.0-ce | 1caf76ce6baa889133ece59fab3c36aaf143d4ef | 341110346f5a7b26d717f13339c6ffeb238e2b00 |
| v17.11.0-ce-rc1 | df7b627909ca07c0a1a000fde37b9a7f264f0e6a | c3544074ccd3149e3a9eef4e909bc9090351e72f |
| v17.11.0-ce-rc2 | d7062e5443166bb01b04d999347d22d742fdf848 | d894a9bc2b0c0bb0577b8471036d115a8d08a56a |
| v17.11.0-ce-rc3 | 5b4af4f7126e185917784eeac63b89b1a03bc9bb | 5f34c1cfcbbd9cfcc1667c85de6e3071471d600c |
| v17.11.0-ce-rc4 | 587f1f003bbe3090710151cc9e89042959d55d0b | 1ce587fbfa0e6262bfbf2b3d96e24160bec39204 |
| v17.12.0-ce | c97c6d62c26c1da407e3086f0b5d3d866ed308bc | 2c27a45733f7217be67a278a6de25c1bf43b489c |
| v17.12.0-ce-rc1 | ee2f9437b6ef584a9db69a910fea90767a030ecb | cb151b253d4b38e2536e94480712396f9f52784f |
| v17.12.0-ce-rc2 | 5d7e8f3778c9dfb65b9dca3cb53a3e5234e9ae84 | 5567d83cb6944652e318d2911f88c2cf7aa7838c |
| v17.12.0-ce-rc3 | 80c8033edeed330cbc2ccdb08aa0b74488816dd4 | cba97ac4bede966e4e4fd54a08b1c006a27a1c87 |
| v17.12.0-ce-rc4 | 6a2c058cd85347b688a4f75089e1c64844966114 | 50f32db3264eb195cc98d7be1bf89454bb46c982 |
| v17.12.1-ce | 7390fc6103da41cf98ae66cfac80fa143268bf60 | 9f395a8a67d4499fb874e20649b92234e26140cc |
| v17.12.1-ce-rc1 | feff709de8555eb0e8c0b4476ed2db8f9c236d79 | fbdca95c6c5cba29e8dc2d5e444f6a2babc3718a |
| v17.12.1-ce-rc2 | b0992d220f60a6591e3293cf8c61bd7eb09ee5ef | c62be9b3b57a7675aaccb922f68b414115a5224b |
| v18.01.0-ce | 03596f51b120095326d2004d676e97228a21014d | dbd165f3d2fd43d80bd4ee19f0d0232015b124e3 |
| v18.01.0-ce-rc1 | 44e2def671d00032daee06eccfff6e6645791d08 | d64390db85eabc78f209c95ca751124912de83e2 |
| v18.02.0-ce | fc4de447b563498eb4da89f56fb858bbe761d91b | c2fc911a77bb061800ceecd9ec7567e54ad8bdbd |
| v18.02.0-ce-rc1 | 5e1d90afd7d77c616c80bada13911a1b96e8c4f2 | 8583d015424c21b3dcdcebb6720af69b87bf9983 |
| v18.02.0-ce-rc2 | f968a2cc556c2c274c23d37518f7f576b8551d03 | a3cc3bef8a49ae6bae274048aa19e66908201701 |
| v18.03.0-ce | 0520e243029d1361649afb0706a1c5d9a1c012b8 | c0c0d291dd210c50f7da581f316d8603107a59d7 |
| v18.03.0-ce-rc1 | 5ff63c0239178c3e6e7e90d592cf5a6af65a7048 | 0b1b7ea8913e7a8f4a2d01f89ed5943198213acf |
| v18.03.0-ce-rc2 | 3e53917a28fc6625a97648c2e7ea3c2f019eb709 | 6de4e3b8b4661d62f6dc98e99228ee7f7af2a801 |
| v18.03.0-ce-rc3 | e7309590a28eaa52fb852d297c2f642c960eb6b1 | 9d277155590e25e4998ae6751be53f2528d245b2 |
| v18.03.0-ce-rc4 | fbedb97a271ba4244334bdd8e1603cda8433c16b | 0a115771172d339dd4624adb92edbabe5c9924eb |
| v18.03.1-ce | 9ee9f402cd1eba817c5591a64f1d770c87c421a4 | 077dedce5df976938a10e453734616627e7cb0f7 |
| v18.03.1-ce-rc1 | dc75023a9a98f1f2b1d2e838577828f8f9597e1a | e9c1dcaa07e811fa85aad8a1fe149bb231294457 |
| v18.03.1-ce-rc2 | fcdc984cfdb8f86df78cbc8853079f06e77cbf95 | a857f1a1155dca76b2cca36beb9175ae404b6cc6 |
| v18.04.0-ce | 3d479c0af67cb9ea43a9cfc1bf2ef097e06a3470 | 9eb621b3483bf820aa3e36c9f00f0781f4054f90 |
| v18.04.0-ce-rc1 | 0c7f7c6ff43a8c9c4d84eb0fb1d9f4f691e2d1ec | b94602d9284126a350c6114fd1396d33b484c2bc |
| v18.04.0-ce-rc2 | f4926a265f6545f1a5e0d200c0a33251d77615b1 | 14fdd61b11740165464b46d78f54415765f9dc2c |
| v18.05.0-ce | 33f00ce1111dcd2dc44b9ab5c71af14b2ce915c5 | 7c86250b8a200f993836aefdb377d2ba7b5237fa |
| v18.05.0-ce-rc1 | 33f00ce1111dcd2dc44b9ab5c71af14b2ce915c5 | 7c86250b8a200f993836aefdb377d2ba7b5237fa |
| v18.09.9-beta1 | 9e36dcdcdfe4e6387d10db5ed3c2a31bb72ecf4a | 1752eb3626e3a17b0881135a9402e57728208fed |

### Packaging

| Tag | Docker-ce sha | docker/docker-ce-packaging sha |
| --- | --- | --- |
| v17.06.0-ce | bb8c1249be4e85c9fd0fc086314452b0b23bb2c4 | d6fd2c20c0cf3d08a5c036233447bbc3c77e1829 |
| v17.06.0-ce-rc1 | 7f8486a39a2a404b70092c8b8767e5ca6dac0889 | f7098d293c50687e63b5c612202bf178710ca96d |
| v17.06.0-ce-rc2 | de513d091d00eb645dd6af7aa907af89dfec0889 | da34e06fd98332e24e38cefdfe1aefb4274c5a1b |
| v17.06.0-ce-rc3 | de513d091d00eb645dd6af7aa907af89dfec0889 | da34e06fd98332e24e38cefdfe1aefb4274c5a1b |
| v17.06.0-ce-rc4 | 29fcd5dfae26aaddbe4f5c84bf913d49db12dde9 | 0bf5b5388e49aaac7c6be48b8c97d8f97a01e0d9 |
| v17.06.0-ce-rc5 | bb8c1249be4e85c9fd0fc086314452b0b23bb2c4 | d6fd2c20c0cf3d08a5c036233447bbc3c77e1829 |
| v17.06.1-ce | 8738f29e5b9ac5cfde0e6a71c917d6d9ee018d4b | 19775056a4898df423921bb5c47870e34e1c997b |
| v17.06.1-ce-rc1 | 77b4dce06699197a5c13b4eb66d7531c00f516c2 | 1183e6389f2e432c23c8beb85c7c6d6e13747d60 |
| v17.06.1-ce-rc2 | 77b4dce06699197a5c13b4eb66d7531c00f516c2 | 1183e6389f2e432c23c8beb85c7c6d6e13747d60 |
| v17.06.1-ce-rc3 | 8738f29e5b9ac5cfde0e6a71c917d6d9ee018d4b | 19775056a4898df423921bb5c47870e34e1c997b |
| v17.06.1-ce-rc4 | 8738f29e5b9ac5cfde0e6a71c917d6d9ee018d4b | 19775056a4898df423921bb5c47870e34e1c997b |
| v17.06.2-ce | 8738f29e5b9ac5cfde0e6a71c917d6d9ee018d4b | 19775056a4898df423921bb5c47870e34e1c997b |
| v17.06.2-ce-rc1 | 8738f29e5b9ac5cfde0e6a71c917d6d9ee018d4b | 19775056a4898df423921bb5c47870e34e1c997b |
| v17.07.0-ce | edadfd04bec1c667a3581be758e3ee88f8924051 | 7c901365c3263747189d6ea8ec8225724a0f68f6 |
| v17.07.0-ce-rc1 | edadfd04bec1c667a3581be758e3ee88f8924051 | 7c901365c3263747189d6ea8ec8225724a0f68f6 |
| v17.07.0-ce-rc2 | edadfd04bec1c667a3581be758e3ee88f8924051 | 7c901365c3263747189d6ea8ec8225724a0f68f6 |
| v17.07.0-ce-rc3 | edadfd04bec1c667a3581be758e3ee88f8924051 | 7c901365c3263747189d6ea8ec8225724a0f68f6 |
| v17.07.0-ce-rc4 | edadfd04bec1c667a3581be758e3ee88f8924051 | 7c901365c3263747189d6ea8ec8225724a0f68f6 |
| v17.09.0-ce | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.09.0-ce-rc1 | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.09.0-ce-rc2 | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.09.0-ce-rc3 | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.09.1-ce | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.09.1-ce-rc1 | 5a2228df09c25512f83542da0f3aa3828ef0a0f1 | cc3be7f52531077a29c2da0deef5f25a8631f3a1 |
| v17.10.0-ce | 73ecdcee1014105418b96dd56ce41b5083019ffc | 2cea55c9434011bf181715fdf6a29e7049fcabfe |
| v17.10.0-ce-rc1 | 5338e5013ab17a326132ebf813fcdc9d7211cbc8 | a1debdbcf0220c10799fce225f2f189e681de385 |
| v17.10.0-ce-rc2 | 73ecdcee1014105418b96dd56ce41b5083019ffc | 2cea55c9434011bf181715fdf6a29e7049fcabfe |
| v17.11.0-ce | 7c931da69ab02ad92040c1828e70834774d3eaa5 | 986aad8023766a04d0cd2f5507f50591a21cbbdd |
| v17.11.0-ce-rc1 | e947e4d4f1a55589a3eb4f049f51ddeddaf8c2da | e2507465495d5361fe9501a5f5e7172779a67e63 |
| v17.11.0-ce-rc2 | 7c931da69ab02ad92040c1828e70834774d3eaa5 | 986aad8023766a04d0cd2f5507f50591a21cbbdd |
| v17.11.0-ce-rc3 | 7c931da69ab02ad92040c1828e70834774d3eaa5 | 986aad8023766a04d0cd2f5507f50591a21cbbdd |
| v17.11.0-ce-rc4 | 7c931da69ab02ad92040c1828e70834774d3eaa5 | 986aad8023766a04d0cd2f5507f50591a21cbbdd |
| v17.12.0-ce | f9cde631a2f2ac685ee6aaf48db2479af0ed8a20 | 13137519814f4b1abbc5e3cc803c04ea4ef40867 |
| v17.12.0-ce-rc1 | cf5224adb217a8a426dba911dc4c3d5d4dffd527 | 6e901d8febcb7c308035952baa760779ab097806 |
| v17.12.0-ce-rc2 | f9cde631a2f2ac685ee6aaf48db2479af0ed8a20 | 13137519814f4b1abbc5e3cc803c04ea4ef40867 |
| v17.12.0-ce-rc3 | f9cde631a2f2ac685ee6aaf48db2479af0ed8a20 | 13137519814f4b1abbc5e3cc803c04ea4ef40867 |
| v17.12.0-ce-rc4 | f9cde631a2f2ac685ee6aaf48db2479af0ed8a20 | 13137519814f4b1abbc5e3cc803c04ea4ef40867 |
| v17.12.1-ce | d70d9c910a4efd1ed60a3a39e0c08f7b940a83dc | c2bf4f23d771878495ef8b7dabb69968fe45fb4c |
| v17.12.1-ce-rc1 | d70d9c910a4efd1ed60a3a39e0c08f7b940a83dc | c2bf4f23d771878495ef8b7dabb69968fe45fb4c |
| v17.12.1-ce-rc2 | d70d9c910a4efd1ed60a3a39e0c08f7b940a83dc | c2bf4f23d771878495ef8b7dabb69968fe45fb4c |
| v18.01.0-ce | 6aa9598694defc590104c2ad86bcf0ac7e726509 | 41ae603d88e02e950ba4f3e227f9f6e972a3c74b |
| v18.01.0-ce-rc1 | 6aa9598694defc590104c2ad86bcf0ac7e726509 | 41ae603d88e02e950ba4f3e227f9f6e972a3c74b |
| v18.02.0-ce | ef4aa7ebe5f5bf8c47ca33dd44a767531304cedc | 7ea33ac7993e2abfd2404e147d95a3b41a29ccbe |
| v18.02.0-ce-rc1 | ef4aa7ebe5f5bf8c47ca33dd44a767531304cedc | 7ea33ac7993e2abfd2404e147d95a3b41a29ccbe |
| v18.02.0-ce-rc2 | ef4aa7ebe5f5bf8c47ca33dd44a767531304cedc | 7ea33ac7993e2abfd2404e147d95a3b41a29ccbe |
| v18.03.0-ce | 0520e243029d1361649afb0706a1c5d9a1c012b8 | 95930e87947055e4d3b8c639f3ee5fc427fb71a0 |
| v18.03.0-ce-rc1 | c160c7335360569036695fa958d327827fcc7dc0 | 138ca8c7ad6e70f522f80444c6508292d3a6fc46 |
| v18.03.0-ce-rc2 | cbc5bef54fee9a1451fd8f224af580c1dc3b28ae | 88176d01f492b112472467a747abc3ba98dd41d8 |
| v18.03.0-ce-rc3 | e7309590a28eaa52fb852d297c2f642c960eb6b1 | 89ec01afcb73f0a7c01816e56cf9cd63a2ed5778 |
| v18.03.0-ce-rc4 | fbedb97a271ba4244334bdd8e1603cda8433c16b | 9cc70ae1b07897d64b48b02b46f6457944fcc3b5 |
| v18.03.1-ce | 9ee9f402cd1eba817c5591a64f1d770c87c421a4 | fdb8850492b82da66ee7f5acecc27731952e8f0d |
| v18.03.1-ce-rc1 | dc75023a9a98f1f2b1d2e838577828f8f9597e1a | 7dd4bb6171adaf3a21022a38a1218779943f35cb |
| v18.03.1-ce-rc2 | fcdc984cfdb8f86df78cbc8853079f06e77cbf95 | b45862f481ecdd86568e01bd15952dc3e8467c5a |
| v18.04.0-ce | 3d479c0af67cb9ea43a9cfc1bf2ef097e06a3470 | 237393c19a148f0ab9d2cc7efc0b549c52e611aa |
| v18.04.0-ce-rc1 | 0c7f7c6ff43a8c9c4d84eb0fb1d9f4f691e2d1ec | 03b5f310ee5eff2a3550524b0b47c46e665c3584 |
| v18.04.0-ce-rc2 | f4926a265f6545f1a5e0d200c0a33251d77615b1 | 456461a7864834883d0abb44d9369d68bb5f75c5 |
| v18.05.0-ce | e58320e94af50e30cf92315f2bed7a1808461b10 | c216602d16cfefdb175cdb902038789b00dd2cef |
| v18.05.0-ce-rc1 | e58320e94af50e30cf92315f2bed7a1808461b10 | c216602d16cfefdb175cdb902038789b00dd2cef |
| v18.09.9-beta1 | 0950f0f1fa79be1405d31ba24bbbe992bc4944a5 | abda7e1227ba45584a6abbf6fbc07e31f970380f |