# ready-git-go
Go tool for updating git repos en masse

## Usage

```
 -P string
        Password for git Auth
  -U string
        Username for git Auth
  -c    Change git repo origin
  -o string
        New origin base file path
  -p string
        Repositories file path
  -q    Quiet output
  -u    Update all repositories in file path (must provide file path)
  -v    Verbose output
  ```

### Chaging the base origin URL for all repos under a base folder

```readygitgo.exe -p c:\Users\dharmaofcode\projects\  -v -c -o "https://newurl.com/Project/_git"```

### Updating all repos under a base folder

```readygitgo.exe -p c:\Users\dharmaofcode\projects\  -v -u```

### Cloning multiple repos at once
TODO
