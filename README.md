# glox
This is a Go implementation of a tree-walk interpreter for the Lox scripting language. 

Lox language specifications and a Java implementation of a tree-walk interpreter for it can be found in the book [Crafting Interpreters](https://craftinginterpreters.com/) by Robert Nystrom.

## Downloading and Installing
To start, you'll need to have Go installed locally.

Using Go run:
```
go get https://github.com/skusel/glox
```

This will download and install the executable to your `$GOPATH/bin`. You can add this directory to your path, if it's not already on it, so you can run `glox` without having to type in the entire path to the executable.

## A Little About Lox
In short, the Lox language is a dynamcially typed, object oriented scripting language with C-like syntax.

When using this interpreter, you'll notice that some things you have come to expect from modern languages and their runtimes are not present. For example the REPL does not remember variables or functions you entered in earlier prompts. There are no built in data structures. Native functions to read and write files do not exist yet. There is no notion of importing code from other source files. With that said, the language has a lot of features built-in already.

## Running the Interpreter
You can run `glox` in two ways.

The first, is via the REPL. To launch the REPL, just type `glox` into your prompt.

The second, is by specifying a `*.lox` file you wish to run.

```
glox /path/to/source.lox
```

The second option, will allow you to dive into the language a lot more. I would recommend using it over the REPL if you are interested in trying this implementation of the language out.

## Lox Examples
This section does not cover all Lox syntax, that's what [Crafting Interpreters](https://craftinginterpreters.com/) is for, but here are some examples of things you can do with the language if you're interested in using this Lox interpreter.

Like most languages Lox allows you to create functions. Here is any example of a Lox function that implements a solution to the well known fizzbuzz question.

```
fun fizzbuzz(target) {
    for(var i = 1; i <= target; i = i + 1) {
        if(i % 3 == 0 and i % 5 == 0) {
            print "fizzbuzz";
        } else if(i % 3 == 0) {
            print "fizz";
        } else if(i % 5 == 0) {
            print "buzz";
        } else {
            print i;
        }
    }
}

fizzbuzz(100);
```

Lox also has many of the object-oriented programming features that will feel familar if you have used other languages like Java, C++, and Python.

```
class Plant {
    init(scientificName) {
        this.scientificName = scientificName;
    }

    howToPlant() {
        print "Dig a small hole in the soil.";
        print "Place plant into the hole.";
    }

    water() {
        print "Once a week.";
    }

    light() {
        print "Requires 4 to 6 hours of sunlight.";
    }
}

class JadePlant < Plant {
    init(scientificName) {
        super.init(scientificName);
        print "Constructing a Jade plant (" + this.scientificName + ")";
    }

    howToPlant() {
        super.howToPlant();
        print "Water lightly.";
    }

    water() {
        print "Every 2 weeks.";
    }
}

var myPlant = JadePlant("Crassula ovata"); // prints "Constructing a JadePlant (Crassula ovata).\n"
myPlant.howToPlant(); // prints "Dig a small hole in the soil.\nPlace plant in the hole.\nWater lightly.\n"
myPlant.water(); // prints "Every 2 weeks.\n"
myPlant.light(); // prints "Requires 4 to 6 hours of sunlight.\n"
print myPlant.scientificName; // prints "Crassula ovata\n"
```

## Structure of the Code
The code structure for this project is relatively flat. `main.go`, which is located in the same directory as this `README.md`, is the entry point to the interpreter. From there you jump into the `lang` directory/package. The Lox source code flows through the scanner, into the parser, then onto the resolver, before being executed in the interpreter. Some other files like `token.go`, `expr.go`, and `stmt.go` are used to represent components of the AST. Logic for callables, native functions, user defined functions, and classes and their instances have also been broken out into their own files. Environments are used to store program state, and they are chained together in a way that reflects the scope of the variables they hold. The `astprinter.go` file was used in earlier stages of development for testing purposes, but is no longer actively used.

## License
Like the Java source code for the original jlox tree-walk interpreter. This glox tree-walk interpreter is also made available under the MIT License. Please see [LICENSE](https://github.com/skusel/glox/blob/main/LICENSE) for more details.