# OpenHDC

OpenHDC (**Open** **H**ybrid **D**ata **C**enter) is an open-source project designed to provide a robust and local server solution for hybrid data sources. 

We preserves data privacy by localizing data and combining siloed information to enhance traceability.

## ✨ Features

- ***Robust Local Server***: Safeguards data with local storage.
- ***Hybrid Data Center***: Consolidates data from multiple sources.
- ***Improved Traceability***: Facilitates data tracking across systems.

## 🍺 Build & Run

1. Install by `go install`

    ```sh
    go install github.com/openhdc/openhdc@latest
    ```

2. Build from source
    1. Clone the repository:

        ```sh
        git clone https://github.com/openhdc/openhdc.git
        ```
    
    2. Change directory and build:

        ```sh
        cd openhdc && make build
        ```
    
    3. Run:
   
        ```sh
        ./bin/openhdc
        ```
    

## 🔨 Environment

Ensure you have the following environment setup:

- Go 1.23.4 or later
- Protobuf compiler (`protoc`)
- Make

## 🔍 Documentation

For detailed documentation, please visit [docs](/docs) directory.

## 🦮 Help

If you need help, feel free to open an issue on GitHub or use the discussions feature to contact the maintainers. 

We'll do our best to assist you promptly.

## 📢 Roadmap
- [ ] v0.0.1
    - [ ] Better error messages
    - [ ] Improved naming conventions
    - [ ] Enhanced application closing procedures

## ⛔ Rules

Please review and adhere to the contribution guidelines outlined in the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## ⚖️ License

This project is licensed under the terms of the [LICENSE](LICENSE) file.


TL;DR

overwrite: upsert
append: insert

incremental: compare with incremental column (default is primary key)


| No     | Sync Mode                 | Primary Key | Incremental Column | Behavior
| :----: |  :----:                   | :----:      | :----:             | :----:
|1       | overwrite                 | v           | no need            | upsert based on primary key, and delete records do not exists
|2       | overwrite                 | x           | no need            | truncate and insert
|3       | append                    | v / x       | no need            | do not touch old records, and insert new record
|4       | incremental_append        | v / x       | x                  | compare with selected incremental column (default is primary key) and insert
|5       | incremental_append_dedupe | v / x       | x                  | compare with selected incremental column (default is primary key) and upsert ???


<!-- |2       | full_append               | v           | x                  | compare with incremental column (default is primary key) and insert, 不管舊資料 -->




| No     | Sync Mode                   | Cursor   | Primary Key | Behavior
| :----: |  :----:                     | :----:   | :----:      | :----:
|1       | overwrite                   |x         | v           | `upsert` based on primary key, and `delete` records do not exists
|2       | overwrite                   |x         | x           | truncate and `insert`
|3       | append                      |x         | v           | do not touch old records, and `insert` new record
|4       | append                      |x         | x           | do not touch old records, and `insert` new record
|5       | append (incremental)        |v         | v           | error
|6       | append (incremental)        |v         | x           | compare with selected incremental column and `insert`
|7       | append_dedupe (incremental) |v         | v           | compare with selected incremental column and `upsert`
|8       | append_dedupe (incremental) |v         | x           | error


| No     | Read Mode   | Write Mode      | Primary Key | Behavior
| :----: | :----:      | :----:          | :----:      | :----: |
|1       | full        | overwrite       | v           | `upsert` based on primary key, and `delete` records do not exists
|2       | full        | overwrite       | x           | truncate and `insert` records
|3       | full        | append          | x           | do not touch old records, and `insert` new records
|4       | incremental | append          | x           | compare with cursor and `insert` records
|5       | incremental | append & dedupe | v           | compare with cursor and `upsert` records








