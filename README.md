# systemctl

## usage

```go
s := systemctl.NewDefault()
s.Start("servicename")
s.Stop("servicename")
s.Enable("servicename")
s.Disable("servicename")

units, err := s.Units()
for _, u := range units {
  // ... handle unit
}

// etc...
```

## TODO

1. Support more .service file fields
2. Read and parse .service files
