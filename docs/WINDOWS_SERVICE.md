# Running Rush as a Windows Service

Rush natively supports running as a Windows Service. Outlined below is are
the general steps to set it up.

1. Obtain the rush windows distribution
2. Create the directory `C:\Program Files\Rush` (if you install in a different
   location simply specify the `-config` parameter with the desired location)
3. Place the rush.exe and the rush.conf config file into `C:\Program Files\Rush`
4. To install the service into the Windows Service Manager, run the following in PowerShell as an administrator (If necessary, you can wrap any spaces in the file paths in double quotes ""):

   ```
   > C:\"Program Files"\Rush\rush.exe --service install
   ```

5. Edit the configuration file to meet your needs
6. To check that it works, run:

   ```
   > C:\"Program Files"\Rush\rush.exe --config C:\"Program Files"\Rush\rush.conf --test
   ```

7. To start collecting data, run:

   ```
   > net start rush
   ```

## Other supported operations

Rush can manage its own service through the --service flag:

| Command                            | Effect                        |
|------------------------------------|-------------------------------|
| `rush.exe --service install`   | Install rush as a service |
| `rush.exe --service uninstall` | Remove the rush service   |
| `rush.exe --service start`     | Start the rush service    |
| `rush.exe --service stop`      | Stop the rush service     |


Troubleshooting  common error #1067

When installing as service in Windows, always double check to specify full path of the config file, otherwise windows service will fail to start

 --config C:\"Program Files"\Rush\rush.conf
