#!/bin/bash
echo Building domain_finder...
go build -o bin/domain_finder

echo Building synonyms...
cd ../synonyms
go build -o ../domain_finder/lib/synonyms

echo Building available...
cd ../available
go build -o ../domain_finder/lib/available

echo Building sprinkle...
cd ../sprinkle
go build -o ../domain_finder/lib/sprinkle

echo Building coolify...
cd ../coolify
go build -o ../domain_finder/lib/coolify

echo Building domainify...
cd ../domainify
go build -o ../domain_finder/lib/domainify

cd ../domain_finder
echo Done.
