### 补充

- git仓库迁移 ： https://blog.csdn.net/zzzgd_666/article/details/81252470
- 回退commit并保存修改 ： git reset --soft commitId
- 回退commit并不保存修改 ： git reset --hard commitId
- 打tag : `git tag -a v1.2 9fceb02`  为某一个commit打tag
	           `git push origin v1.5`	推送这个tag
			   删除tag
- git 外部仓库：https://www.yiibai.com/git/git_submodule.html


### git命令复习

branch

- 删除本地分支：git branch -D BranchName
- 删除远程分支：git push origin --delete BranchName
- 建立本地分支：git checkout -b BranchName (-b 新建分支，并切换分支)
- 推送到远程：git push origin LocalBranchName:RemoteBranchName
- 关联本地分支到远程分支：git branch --set-upstream-to=origin/RemoteBranchName


更新本地的代码为最新的远程分支代码（在当前分支下）

- git fetch origin RemoteBranchName:LocalBranchNameTemp
- git diff LocalBranchNameTemp
- git merge LocalBranchNameTemp
- git branch -D LocalBranchNameTemp

在本地合并两个远程分支

- 切换到主分支 git checkout RemoteBranchName1 会建立一个本地分支
- 切换到副分支 git checkout RemoteBranchName2 会建立一个本地分支
- 切换到主分支
- merge副分支
- vscode左侧边栏 提示merge的冲突信息包括（新文件和已有文件的修改信息）其中新文件可以不用管，但是已有文件的修改需要手动解决冲突。首先讲冲突文件从staged changes移出到changes，在changes中对照修改冲突（staged changes不可修改）。然后填写vscode下边栏中填写merge信息，并merge。最后将已有冲突文件push到远程分支。 
- push到远程					


版本回退

- 本地回退 git reset --hard commitID
- 远程回退 git reset --hard commitID
                    git push origin LocalBranch:RemoteBranch --force


查看本地分支对应的远程分支 git branch -vv

commit 之后撤回到add状态 git reset --soft HEAD^ 




tag 与 branch

- tag是静态的，为当前代码做标记，可以根据tag方便地回滚到tag标记的位置。tag是某次commit的指针。创建tag、推送tag、删除tag： [tag基本操作](https://blog.csdn.net/beyond702/article/details/78304326) 
- branch是动态的，在不断的更新开发，用于多人开发后merge到master。git衍和可以完美解决正在别的分支开发时，matser更新的问题。[git详细讲解](https://www.cnblogs.com/guge-94/p/11281724.html) 


fork和clone

- 在没有被授权时，clone下来的代码无法通过git push提交修改。此时可以通过fork后，再clone自己fork的代码，进行修改后，git push到自己的仓库，在使用pull request命令，将自己的修改申请合并到别人的代码中。


issue用于讨论