# Test Acq Dates

![technology Go](https://img.shields.io/badge/technology-go-blue.svg)

This is an application to test [acq-dates](https://github.com/mercadolibre/fury_acq-dates).

## Input File

For this version, it is necessary that the input file should be named paths.csv.

This file must be a CSV file with two fields:

- ARN: Item's ARN.
- URL: acq-dates' endpoint that needs to be tested.

Currently, URL is generated through Google Spreadsheets.

## How to use?

1 - Clone the repository

```sh
git clone https://github.com/tales-lopes-meli/testing-acq-dates.git
```

2 - Generate paths.csv

- It is necessary to store paths.csv in root directory.

3 - Execute the script

```sh
sh start.sh
```

4 - Check out your results

## Attention

It necessary to enable VPN to execute correctly.
