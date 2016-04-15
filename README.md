# checker

Let’s assume you have hundreds of thousands of servers under your command and
all servers must be within the following standards:

- Given some files, must be in defined paths
- Given some strings, must be in given path's content (eg: in a log file)
- Given some process names, must be running

So checker is a simple kite with a check method, making given checks of the
sort:

```javascript
{
  "check_etc_hosts_has_8888": {
    "path": "/etc/hosts",
    "type": "file_contains",
    "substring": "8.8.8.8"
  },
  "check_kite_config_file_exists": {
    "path": "/etc/host/koding/kite.conf",
    "type": "file_exists"
  },
  "check_go_is_running": {
    "name": "go",
    "type": "process_is_running"
  }
}
```

## Example usage

Suppose we know the host and port of a checker kite. We can run checks agains
it with something like:

```go
func check(host string) {
	k := kite.New("investigator", "1.0.0")

	investigateWorker := k.NewClient(host)
	investigateWorker.Dial()

	checks := map[string]handlers.CheckRequest{
		"go_is_running": {
			Name: "go",
			Type: "process_is_running",
		},
		"checks_native_go_exists": {
			Path: "./checks/native.go",
			Type: "file_exists",
		},
		"checks_native_go_has_func": {
			Path:   "./checks/native.go",
			Substr: "func",
			Type:   "file_contains",
		},
		"checks_missing_file": {
			Path: "/tmp/such_hax0r_shell",
			Type: "file_exists",
		},
	}

	response, err := investigateWorker.Tell("local_check", "native", checks)
	if err != nil {
		k.Log.Fatal("err: %s\n", err)
	}

	k.Log.Info("response: %s\n", response)
}
```


Expected result from such a call would be:

```javascript
{
  "checks_missing_file": {
    "result": false,
    "error": {
      "Op": "stat",
      "Path": "\/tmp\/such_hax0r_shell",
      "Err": 2
    }
  },
  "checks_native_go_exists": {
    "result": true,
    "error": null
  },
  "checks_native_go_has_func": {
    "result": true,
    "error": null
  },
  "go_is_running": {
    "result": true,
    "error": null
  }
}
```

Using a function like this example `check()` (except we might want to have
these checks parametrized) would be easy to be used in a fan-in -> fan-out
manner.

## Future improvements

- Add a `RemoteCheck()` handler so one kite would be able to check other kites,
  in order to balance the IO load between more machines, without changing the
  fact that only one machine is enough to initiate the whole thing.

- Add a feature for changing the direction by heart-beating check results to a
  broker once every `N` minutes, instead of waiting for request.
