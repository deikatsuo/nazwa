/* Nazwa
*  Membuat tabel utama
*  Authored by Deri Herdianto
*  c 2020
*/

-- name: create-role-table
CREATE TABLE `role` (
    `id` SMALLINT GENERATED ALWAYS AS IDENTITY,
    `name` VARCHAR(10),
    PRIMARY KEY (id)
);