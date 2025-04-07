# Data Types

## Organization

* id - primary key int
* name - string
* owners - []user (maybe its own table)
* locations - []location

## Location

* id - primary key int
* organization - organization fkey
* address - string
* city - string
* state - string

## User

* id - primary key
* firstName - string
* lastName - string

## Auth/Account

email - string
password - string
salt - string
active

## Vendor

* id
* name
* ...

## Organization Vendor

* organizationId
* vendorId

## Brand

* id
* name

## OrgTag

* id
* name
* org - Organization
* type - category that tag applies to (ingredient, product, ...)

## Unit (does this belong in the db??)

* id
* name
* type - (mass, volume, ...(case, pack)) // is this necessary?

## Ingredient

If i order a case of cups, thats 1000 cups, so the cup ingredient would be 1000 cups for $n.
10 tomatoes
1 bundle of bananas
1000 grams of flour

* id
* org
* name
* brand - Org Brand
* vendor - Org Vendor
* category  - Tag
* unitType - Unit
* unitCount - number // how many units purchased
* purchasePrice

## Product

* id
* org
* name
* menuPrice
* servings

## Product Ingredient

* id
* product
* ingredient
* amount
