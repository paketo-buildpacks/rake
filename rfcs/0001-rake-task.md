{{Runtime Raketask Specification}}
Proposal
Allow app users to specify a raketask from the rakefile at runtime.

Motivation
Specifying an argument to `bundle exec rake` is a common use-case. Allowing the task to be
specified at run time makes the container image more flexible/broadly useable because it enables
the non-default use-case.

Implementation (Optional)
The user can specify the rake task as an arguement to the launcher at runtime.
We believe [this rfc](https://github.com/buildpacks/rfcs/blob/main/text/0045-launcher-arguments.md) will enable this behavior.

For example, if the user builds a `rake` buildpack with a Rakefile that has the tasks `task1` and `task2`,
they can run task1 with:

```bash
$ docker run image-name task1
```

task2 with
```bash
$ docker run image-name task2
```

and simply running
```bash
$ docker run image-name
```
will run whatever is defined as the default task.

