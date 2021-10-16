
# Models package

## IDType

`IDType` is a simple type alias that currently equates to `int64`.  It might be removed in the future, but for now allows us to adjust the type used for model instance IDs with a single-line change, rather than a sweeping edit of the entire codebase.

## SelectQuery{}

`SelectQuery{}` is a helper struct that encapsulates the parameters needed to execute a SELECT on a given table.

In general, when a model method fetches a list of rows (like the `.All` family of methods), that method requires the caller to pass in a `SelectQuery{}`, usually containing `Limit` and `Offset` parameters extracted from the HTTP request.  Additionally, some API calls may set the `OrderBy` and `Descending` parameters as well.

If a model method only fetches a single row (i.e. the `.Get` family of methods), it will not require a `SelectQuery{}` struct to be passed in.  In this case, the struct is used internally by the method, both for consistency's sake as well as to avoid bugs resulting from writing queries manually.

To set the main body of the query (i.e., `SELECT * FROM table WHERE...`), call `.SetQuery("SELECT * FROM table WHERE...", arg1, arg2, arg3)` with the appropriate args.  The variable substitions use the typical dollar-prefixed syntax (`$1`, `$2`, etc.).

Finally, to retrieve the query's full text (for passing into a `db.Get` or `db.Select` call), call `SelectQuery.Query()` and `SelectQuery.Args()` like so:

```go
err := db.Select(&persons, query.Query(), query.Args()...)
```

The query text will be constructed (included `LIMIT`, `OFFSET`, `ORDER BY`, and `DESCENDING` clauses, if appropriate), and the array of args will be returned in the proper order.

## Tests

The model tests are set up so that they test two things:

1) That each model method makes the expected SQL queries (with the expected arguments)
2) That each model method returns the expected values

They are set up in a typical TDD format, using the Ginkgo package.  Each expectation is written out in plain English.  Each test is run independently (meaning setup and teardown happen for every single expectation to ensure a clean environment).

### General notes

Here are some things to keep in mind that hopefully will help with writing, debugging, and refactoring tests:

1) **SQL expectations are extremely rigid.**  The queries have to occur in *exactly* the order they're expected.  All of them *must* occur.  Any unexpected queries *must not* occur.  The arguments must match *exactly*, etc.  When writing or debugging these, it's useful to have the model method and the test open side-by-side to make sure that the test actually replicates the exact queries made by the model, and does so in exactly the same order.

2) **Many model methods call other model methods** (for example, `.Get`ting an object may entail `.Get`ting its subobjects as well).  For example, at the time of this writing, `MusicGroup.All` calls `MusicGroup.GetImage`.  So does `MusicGroup.Get`.  So does `MusicGroup.AllByPersonID`.  In order to avoid rewriting the `GetImage` query expectation for every one of these functions, just make a helper method that sets up the expectation for you.  In `musicgroup_test.go` you'll see this:

    ```go
    func expectQuery_MusicGroup_GetImage(mock sqlmock.Sqlmock, imageobjectID IDType) {
        mock.ExpectQuery(`SELECT (.+) FROM imageobject`).
            WithArgs(imageobjectID).
            WillReturnRows(
                mockResultRows(
                    imageobjects[imageobjectID],
                ),
            )
    }
    ```

    This function is reused many times in the tests to avoid duplication.

3) **Most of the `time.Time` values in the tests are mocked using `time.Now()`.**  Because of thermodynamics, it's easier to expect these values using the `AnyTime{}` matcher (which is defined in `index_test.go`).  An example of using this matcher would be:

    ```go
        mock.ExpectExec(`UPDATE imageobject`).
            WithArgs(mg.Image.CID, mg.Type, mg.Context, AnyTime{}, mg.Image.ContentURL, mg.Image.EncodingFormat, mg.Image.ID).
            WillReturnResult(sqlmock.NewResult(123, 1))
    ```

    We could also test this in other ways, for example, picking a specific time/date.

4) **Every SQL query returns something, and in tests, you have to decide what it is.**  When you do so, it causes the `db.Query` or `db.Exec` functions to return whatever you specify as the return value.

    For queries that are executed with the `db.Query` method, the return value should be one or more rows, which are specified by using the `mockResultRows` helper method along with one or more fixtures:

    ```go
    mock.ExpectQuery(`SELECT (.+) FROM musicgroup`).
        WithArgs(musicgroupID).
        WillReturnRows(
            mockResultRows(
                musicgroups[musicgroupID],
            ),
        )
    ```

    When the above code runs, it actually causes the query to return the `musicgroups[musicgroupID]` fixture to the model, which then converts it to a struct and returns it to the test.

    For `db.Exec` queries, you'll do something slightly different, as they don't return multiple rows:

    ```go
    mock.ExpectExec(`INSERT INTO musicgroup_members`).
        WithArgs(member.ID, member.Description, member.PercentageShares, member.MusicGroupAdmin).
        WillReturnResult(sqlmock.NewResult(123, 1))
    ```

    This query returns a single integer (representing the ID of the new row).  Generally speaking, I've found that the values you pass to `WillReturnResult` don't matter at all, and are never tested later.  I tend to use `sqlmock.NewResult(123, 1)` for almost all of them.

5) **All of the fixtures used in the tests are defined in `fixtures_test.go`.**  I've strived to keep these internally consistent (meaning that IDs of subobjects actually point to other fixtures with those IDs, etc.).  Why?  Because sometimes, models with subobjects will use these IDs to fetch those subobjects.  In the tests, we want to make sure the models are grabbing the right stuff.  So our fixtures have to be as realistic as a production database, at least in certain ways.

    When you see values in the fixtures like `"_family_name"` or `"_catalog_number"` (in other words, an underscore-prefixed version of the field name), those are simply dummy data.  I prefixed them with underscores to make it clear to anyone working on the tests that they're looking at the value, not the key.

### Test errors

Most test errors are fairly easy to read and interpret.  However, the package we use to test SQL queries generates errors that are very verbose.  **If you get a test error with multiple screens worth of output, it is almost certainly from a failed SQL expectation.**  The multiple screens of output result from the fact that these errors will actually list out the entire contents of the objects (in pretty-printed JSON) that they're dealing with ... twice!  Once with the expected object, and once for the received object.  You can pick through this JSON to see if there are fields containing incorrect values, but in my experience thus far, it's much easier to try other things first and to save this as a last resort.

Failed SQL expectations tend to arise from a few sources:
- the expected query did occur, but it was off by a character or two (they're specified using case-sensitive regex, so this is a good place to start)
- one of the arguments was missing or wrong, or there was an extra argument
- the expected query didn't occur (check with the model method open side-by-side to make sure you're expecting everything properly)
- you used `.ExpectQuery` instead of `.ExpectExec`, or vice versa (you have to check the model method to see which one is being used... unfortunately it's not standardized in our codebase at the moment)




# Authors/contributors

If the above is incomprehensible, reach out to one of the following folks:

- Bryn Bellomy (<bryn.bellomy@gmail.com>)

