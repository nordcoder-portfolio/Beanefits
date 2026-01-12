-- +goose Up
-- +goose StatementBegin
DO
$$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
            CREATE TYPE event_type AS ENUM ('EARN', 'SPEND');
        END IF;
    END
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO
$$
    BEGIN
        IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'event_type') THEN
            DROP TYPE event_type;
        END IF;
    END
$$;
-- +goose StatementEnd
