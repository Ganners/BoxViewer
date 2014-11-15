Box Viewer implementation in Go to replace Google Docs
======================================================

This has been built as a simple client of Box View which should accept a query
string and become something that can be directly embedded as an iFrame.

The purpose is to replace Google Docs viewer, which only has a limited number of
views before it shows a bandwidth error - not good! We found that we had to
implement Box View individually in a multitude of applications which wasn't
viable in the long run.

Requires no external dependencies. To run it just supply an API key and
a location to store the downloaded files (and the data) which has
permission:

    $ boxviewer --key="yourkeyhere" --location="files"
