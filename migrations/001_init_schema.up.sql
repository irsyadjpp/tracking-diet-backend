-- ===========================================================
-- 001_init_schema.up.sql
-- Tracking Diet Backend - Initial Schema
-- ===========================================================

-- =========================================
-- EXTENSIONS (opsional, tapi berguna)
-- =========================================
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";  -- untuk UUID jika dibutuhkan
CREATE EXTENSION IF NOT EXISTS "pgcrypto";   -- untuk hashing password dll

-- =========================================
-- 1. USERS TABLE
-- =========================================
CREATE TABLE users (
    user_id         SERIAL PRIMARY KEY,
    full_name       VARCHAR(100) NOT NULL,
    email           VARCHAR(120) UNIQUE NOT NULL,
    password_hash   TEXT NOT NULL,
    gender          CHAR(1) CHECK (gender IN ('M','F')),
    birth_date      DATE,
    height_cm       DECIMAL(5,2),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =========================================
-- 2. BODY MEASUREMENTS
-- =========================================
CREATE TABLE measurements_body (
    body_id         SERIAL PRIMARY KEY,
    user_id         INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    measured_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    weight_kg       DECIMAL(5,2),
    bmi             DECIMAL(4,2),
    waist_cm        DECIMAL(5,2),
    hip_cm          DECIMAL(5,2),
    arm_cm          DECIMAL(5,2),
    thigh_cm        DECIMAL(5,2),
    body_fat_pct    DECIMAL(4,1),
    muscle_mass_kg  DECIMAL(5,2),
    visceral_fat    DECIMAL(4,1),
    whr             DECIMAL(4,2),           -- waist-hip ratio
    skinfold_mm     DECIMAL(5,2)
);

-- =========================================
-- 3. NUTRITION / ENERGY
-- =========================================
CREATE TABLE measurements_nutrition (
    nutrition_id        SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    measured_at         DATE NOT NULL,
    calories_in_kcal    INT,
    carbs_g             DECIMAL(6,1),
    protein_g           DECIMAL(6,1),
    fat_g               DECIMAL(6,1),
    fiber_g             DECIMAL(6,1),
    water_intake_l      DECIMAL(4,2),
    calories_out_kcal   INT,
    step_count          INT,
    meal_timing_note    TEXT
);

-- =========================================
-- 4. METABOLIC HEALTH
-- =========================================
CREATE TABLE measurements_metabolic (
    metabolic_id        SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    measured_at         DATE NOT NULL,
    systolic_bp         SMALLINT,
    diastolic_bp        SMALLINT,
    fasting_glucose     DECIMAL(5,2),
    hba1c_pct           DECIMAL(4,2),
    cholesterol_total   DECIMAL(5,2),
    ldl                 DECIMAL(5,2),
    hdl                 DECIMAL(5,2),
    triglycerides       DECIMAL(5,2),
    uric_acid           DECIMAL(4,2),
    liver_function_note TEXT,
    kidney_function_note TEXT
);

-- =========================================
-- 5. FITNESS & PERFORMANCE
-- =========================================
CREATE TABLE measurements_fitness (
    fitness_id          SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    measured_at         DATE NOT NULL,
    vo2max_ml_kg_min    DECIMAL(5,2),
    strength_1rm_kg     DECIMAL(5,2),
    pushup_count        INT,
    squat_count         INT,
    endurance_note      TEXT,
    flexibility_note    TEXT
);

-- =========================================
-- 6. WELLBEING / SUBJECTIVE
-- =========================================
CREATE TABLE measurements_wellbeing (
    wellbeing_id        SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    measured_at         DATE NOT NULL,
    energy_level        SMALLINT CHECK (energy_level BETWEEN 1 AND 10),
    mood_level          SMALLINT CHECK (mood_level BETWEEN 1 AND 10),
    hunger_level        SMALLINT CHECK (hunger_level BETWEEN 1 AND 10),
    sleep_hours         DECIMAL(4,2),
    sleep_quality       SMALLINT CHECK (sleep_quality BETWEEN 1 AND 10),
    digestion_note      TEXT,
    stress_level        SMALLINT CHECK (stress_level BETWEEN 1 AND 10)
);

-- =========================================
-- 7. LAB TESTS (Optional)
-- =========================================
CREATE TABLE lab_tests (
    lab_id              SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    test_name           VARCHAR(100) NOT NULL,
    result_value        VARCHAR(50),
    unit                VARCHAR(20),
    reference_range     VARCHAR(50),
    measured_at         DATE NOT NULL
);

-- =========================================
-- 8. AI RECOMMENDATIONS (Genkit integration)
-- =========================================
CREATE TABLE ai_recommendations (
    id                  SERIAL PRIMARY KEY,
    user_id             INT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    generated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    category            VARCHAR(50),      -- e.g. "meal_plan", "progress_report"
    input_summary       JSONB,
    output_text         TEXT,
    confidence_score    DECIMAL(3,2)
);

-- =========================================
-- INDEXES (opsional untuk performa)
-- =========================================
CREATE INDEX idx_body_user_measured_at
    ON measurements_body (user_id, measured_at DESC);

CREATE INDEX idx_nutrition_user_measured_at
    ON measurements_nutrition (user_id, measured_at DESC);

CREATE INDEX idx_metabolic_user_measured_at
    ON measurements_metabolic (user_id, measured_at DESC);

CREATE INDEX idx_wellbeing_user_measured_at
    ON measurements_wellbeing (user_id, measured_at DESC);
