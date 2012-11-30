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
or specify it when running your code:(http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorsForZip)

    $ SUNLIGHT_KEY=yourKey go run hellosunlight.go


### Getting Legislators

#### Searching for Legislators

You can specify any of the fields on the
[Legislator](http://go.pkgdoc.org/github.com/adharris/gosunlight#Legislator) type, and
use the [LegislatorGetList](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetList)
function.  So to get all Democratic members of congress from New York:

    toMatch := gosunlight.Legislator{Party:"D", State: "NY"}
    legislators, err := gosunlight.LegislatorGetList(toMatch)

If it is known that there will only be one result (i.e. when matching to an ID), you
can use the [LegislatorGet](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGet)
function.

    toMatch := gosunlight.Legislator{BioguideID: "S000148"}
    legislator, err := gosunlight.LegislatorGet(toMatch)

However, this function will error if more than one legislator matches:

    toMatch := gosunlight.Legislator{FirstName: "John"}
    legislator, err := gosunlight.LegislatorGet(toMatch)
    if err != nil {
      // There are lots of Johns in DC!
    }

By default, LegislatorGet and LegislatorGetList return only members of the
current congress. Sunlight does have data on past congresses, which can be
accessed by using [LegislatorGetListAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetListAll)
and [LegislatorGetAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorGetAll).
These functions behave the same as above, but include all legislators in
Sunlight's database.

There are two convenience functions provided for searching for legislators:
[Legislator.Get](http://go.pkgdoc.org/github.com/adharris/gosunlight#Legislator.Get) and
[Legislator.GetAll](http://go.pkgdoc.org/github.com/adharris/gosunlight#Legislator.GetAll).
These functions will use the LegislatorGet(All) functions to load a single
legislator in place:

    nancy := gosunlight.Legislator{FirstName:"Nancy", LastName:"Pelosi"}
    nancy.Get() // populate the rest of the fields
    fmt.PrintLn(nancy) // Rep Nancy Pelosi (D CA)

#### Fuzzy Search

Sunlight provides fuzzy searching on Legislator Names.  This can be done with
the [LegislatorSearch](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorSearch)
function:

    legislators, err := gosunlight.LegislatorSearch("Reed")
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

#### Searching by Zip Code

You can get all Legislators for zip code using the
[LegislatorsForZip](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorsForZip)
function.  This returns all legislators for a five digit zip code.

Because 5 digit zips do not uniquely determine a congressional district, it is
possible for this function to return multiple members of the House.  Each House
Member that represents at least part of a zip code will be returned

    legislators, err := gosunlight.LegislatorsForZip("94121")

#### Searching by Latitude and Longitude

You can get precise legislator information if you geocode the address and
search by latitude/longitude using the
[LegislatorsForLatLong](http://go.pkgdoc.org/github.com/adharris/gosunlight#LegislatorsForLatLong)
function.  Sunlight will match the coordinates against the districts, and
return legislators for that district.  This typically means 2 Senators and
1 Representative.

    legislators, err := gosunlight.LegislatorsForLatLong(35.778788, -78.787805)

### Districts

#### Districts by Zip Code

You can get a list of districts that contain any part of a zip code using
the [DistrictsFromZip](http://go.pkgdoc.org/github.com/adharris/gosunlight#DistrictsFromZip)
function:

    districts, err := gosunlight.DistrictsFromZip("12345")

#### Districts by Latitude/Longitude

You can get the single district that contains a point using the
[DistrictsFromLatLong](http://go.pkgdoc.org/github.com/adharris/gosunlight#DistrictsFromLatLong)
function. This will return the district that contains the point for the current
congress:

    district, err := gosunlight.DistrictsFromLatLong(35.778788,-78.787805)

For the 2012 election, the districts were redrawn. Until the new congress is
sworn in in January 2013, DistrictsFromLatLong will return the districts as
they were for the 2010 election.  To get the new districts in the mean time,
[DistrictsFromLatLong2012](http://go.pkgdoc.org/github.com/adharris/gosunlight#DistrictsFromLatLong2012)
will return the new districting.

### Committees

#### Listing Committees

Committees are for one of three chambers: House, Senate, or Joint. A List of
all committees and their subcommittees can be retrieved using the
[CommitteeGetList](http://go.pkgdoc.org/github.com/adharris/gosunlight#CommitteeGetList)
function.

    committees, err := gosunlight.CommitteeGetList("House")

Note that even though the gosunlight Committee type has a field for members
of a committee, that field is *not* populated by this function.  You can
populate the empty members field using committee.GetMembers() function.

#### Getting a Committee

To load a specific committee, its subcommittees, and its members, you can use
the [CommitteeGet](http://go.pkgdoc.org/github.com/adharris/gosunlight#CommitteeGet)
function.  This function takes a committee Id as a parameter:

    committee, err := gosunlight.CommitteeGet("JSEC")

### Getting committees by legislator

To get all the committees and subcommittees a legislator serves on, use the
[CommitteesForLegislator](http://go.pkgdoc.org/github.com/adharris/gosunlight#CommitteesForLegislator)
function.  This function needs the Bioguide ID for the legislator:

    committees, err := gosunlight.CommitteesForLegislator("S000148")

Additionally, if you have a Legislator object from one of the legislator
functions, you can fetch the committees that legislator serves on using the
Committees method:

    toMatch := gosunlight.Legislator{FirstName:"Nancy", LastName:"Pelosi"}
    nancy, _ := gosunlight.LegislatorGet(toMatch)
    committees, err := nancy.Committees()

When committees are fetched this way, the list of committees is cached in the
legislator object, so subsequent calls to Committees() will not result in
additional requests to Sunlight.