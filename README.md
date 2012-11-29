gosunlight - Sunlight Labs Congress API in Go
==========

The [Sunlight Labs Congress API](http://services.sunlightlabs.com/docs/Sunlight_Congress_API/)
provides basic information about Members of Congress. (e.g. party, district,
various websites, etc...)

Use of the Sunlight Labs API requires an API key, which can be acquired
for free by signing up here](http://services.sunlightlabs.com/docs/Sunlight_Congress_API/)

Godoc is available at http://go.pkgdoc.org/github.com/adharris/gosunlight

### Installing go sunlight

Simply fetch the package using `go get`:

    $ go get github.com/adharris/gosunlight


### Setting the Sunlight API Key

After [signing up for an API Key](http://services.sunlightlabs.com/docs/Sunlight_Congress_API/),
you must set it in gosunlight.  This can be done one of two ways:

#### Set the SunlightKey package variable

This can be done almost anywhere, as long as it is done before you make your
first call to any gosunlight function.

    package main

    import "github.com/adharris/gosunlight"

    func init() {
      gosunlight.SunlightKey = "your api key"
    }

#### Set a SUNLIGHT_KEY environment variable

If no key is specificed, gosunlight will look for an environment variable
named `SUNLIGHT_KEY`.  You can either set this in a .bashrc or equivalent,
or specify it when running your code:

    $ SUNLIGHT_KEY=yourKey go run hellosunlight.go


### Getting Legislators

#### Searching for Legislators

You can specify any of the fields on the
[Legislator](http://go.pkgdoc.org/github.com/adharris/gosunlight#Legislator) type, and
use the [LegislatorGetList](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetList)
function.  So to get all Democratic members of congress from New York:

    toMatch := Legislator{Party:"R", State: "NY"}
    legislators, err := LegislatorGetList(toMatch)

If it is known that there will only be one result (i.e. when matching to an ID), you
can use the [LegislatorGet](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGet)
function.

    toMatch := Legislator{BioguideID: "S000148"}
    legislator, err := LegislatorGet(toMatch)

However, this function will error if more than one legislator matches:

    toMatch := Legislator{FirstName: "John"}
    legislator, err := LegislatorGet(toMatch)
    if err != nil {
      // There are lots of Johns in DC!
    }

By default, LegislatorGet and LegislatorGetList return only members of the
current congress. Sunlight does have data on past congresses, which can be
accessed by using [LegislatorGetListAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetListAll)
and [LegislatorGetAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetAll).
These functions behave the same as above, but include all legislators in
Sunlight's database.

#### Fuzzy Search

Sunlight provides fuzzy searching on Legislator Names.  This can be done with
the [LegislatorSearch](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorSearch)
function:

    legislators, err := LegislatorSearch("Reed")
    // will return both senators Reed and Reid

Just like the above, use [LegislatorSearchAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorSearchAll)
to search for legislators in past congresses.

Sunlight fuzzy search also allows for a threshold to be specified, which is
a number between 0 and 1 with 1 being a "perfect match."  Gosunlight defaults
the threshold to .8 (which is also the Sunlight default), but this can be
overridden by changing the package level variable LegislatorSearchTheshold:

    gosunlight.LegislatorSearchTheshold = .9
    // legislator fuzzy search will now be more strict

Sunlight does not recommend values less than .8
