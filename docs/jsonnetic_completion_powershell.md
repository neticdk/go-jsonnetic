## jsonnetic completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	jsonnetic completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
jsonnetic completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -d, --debug               Debug mode
      --log-format string   Log format (plain|json) (default "plain")
      --log-level string    Log level (debug|info|warn|error) (default "info")
      --no-color            Do not print color
```

### SEE ALSO

* [jsonnetic completion](jsonnetic_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 18-Mar-2025
