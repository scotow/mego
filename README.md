# mego 💾

Mego is a simple [megatools](https://megatools.megous.com) command wrapper, allowing you to use the `megatools dl` command with a download list of links and add an auto-try tool.

### Ideas

Megatools is a collection of programs for accessing Mega.nz service from a command line of your desktop or server.

#### Auto-retry

While using the `megatools dl` command to download a bunch of large files, I often found myself being blocked by Mega because I exceeded the bandwidth limit (aka. error 509).

Indeed, Mega allows users to download a few (apparently not fixed) number of GB per day (once again, apparently not fixed).

By default `megatools dl` only retries 3 times when this error occurred, preventing the download of file during the night or while being away from the computer. To fix this problem, `mego` check the error code returned by `megatools dl`, and retry if the command failed.

#### List of files

Another problem that I found while using `megatools dl` is the lack of options to download multiple files at once, and keeping track of the done ones.

To solve this problem, `mego` accepts as arguments a path of file(s) containing a list of Mega download links. `mego` will open the file and start downloading the files listed in it. Once the download of the file successfully terminated, `mego` will add a `#` before the link and write it in the file, preventing the next execution to re-download the file. `mego` will also mark invalid links it found with the `#-` string.

#### Compatibility

Because this script is a wrapper around the `megatools` command, it heavily depends on the outputs of the command. If you have problems using this script, be sure to use the version 1.11.0 (04.04.20) of `megatools`. You can download the latest version [here](https://megatools.megous.com/builds/experimental/).

### Usage

```
Usage of mego:

  mego [-c COMMAND_PATH] [-s SPEED] [-p] [-r INTERVAL] LINK... LINK_PATH...

Application Options:
  -s, --speed-limit=SPEED       Speed limit passed to megatools dl as --limit-speed (default: 0)
  -p, --pipe-outputs            Pipe megatools's stdout and stderr
  -r, --retry=INTERVAL          Interval between two retries (default: 15min)
  -c, --command=COMMAND_PATH    Path to the megatools command (default: megatools)
```

NB: The whole content of a *list file* is read and kept in memory. Every time a file is downloaded, the content of the *list file* will be overwritten. So please do not use a *list file* as a queue during execution.  
