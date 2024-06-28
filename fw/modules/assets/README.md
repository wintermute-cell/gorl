# Module: assets

The `assets` module is a wrapper around file loading operations such as
`os.Open()` or `rl.LoadTexture()`. These wrappers are usually drop in
replacements for the original function.

Using the `assets` module allows the usage of packfiles. When
`assets.UsePackfile()` is called at the start of the program, the program
expects a `data.pack` file to be present in the build output, and will try to
read all file data from that packfile.

A packfile can be build with the tool, using `tool packer <in_dir> <out_file>`,
where `<in_dir>` is usually the `assets` directory and `<out_file>` should be
`build/data.pack`.
