# DroneVis

When working with a large monorepo with several build artifacts
that uses a heaping spoonful of `depends_on` within [drone](github.com/drone/drone)
it becomes hard to catch errors when writing a pipeline.

`DroneVis` will render your pipeline to an image graphviz style.

# Installation

```  
go get github.com/adityanatraj/dronevis
```

# Building the binary

```
go build -o dronevis cmd/main.go
```

# Usage

```
cat my_drone_pipeline.yml | ./dronevis
```

or 

```
./dronevis <path-to-my-pipeline.yml>
```

# Example

You can see an example pipeline in `example/pipeline.yml` that was
rendered to the pngs for each drone condition-pathway.

_note_: much of the content of the steps were removed as this is a real life
pipeline with (perhaps) sensitive details. The fields left out are irrelevant 
to processing as of today.

# Contributing

This project is extremely alpha and could use help. 
Any and all contributions via pull request are welcome!

Possible todos:
- be any amount smarter about step when conditions
- allow specifying where the files are saved
- allow specifying what the output files are named
- allow specifying the filetype of the output image
- allow specifying the graphviz layout engine used
- allow outputting the "dot" file directly to pipe into graphviz
- visualize "failure" transitions
- maybe use svg with tooltips for details

# Thanks

this project wouldn't be possible without:
- [DroneCI](github.com/drone/drone)
- [go-graphviz](https://github.com/goccy/go-graphviz)
- [graphviz](graphviz.gitlab.io)

# License

tl;dr MIT License. 2020 Aditya Natraj

Please see the attached `LICENSE` file.
