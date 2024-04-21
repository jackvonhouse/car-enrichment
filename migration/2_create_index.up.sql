BEGIN;

CREATE INDEX IF NOT EXISTS idx_car_regnum ON car (regNum);

COMMIT;
