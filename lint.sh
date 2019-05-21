errcheck .
go vet .
golint .
revive -config revive.toml .
gosec .
drygopher
