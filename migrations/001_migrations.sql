-- Единая схема для всех сервисов AutoJong
CREATE SCHEMA IF NOT EXISTS autojong;

-- Создание таблицы requests, если она не существует
CREATE TABLE IF NOT EXISTS requests (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL,
    email TEXT NULL,
    car_info TEXT NULL,
    date TIMESTAMP NOT NULL
);

-- Создание индексов
CREATE INDEX IF NOT EXISTS idx_requests_id ON requests(id);
CREATE INDEX IF NOT EXISTS idx_requests_phone ON requests(phone);
