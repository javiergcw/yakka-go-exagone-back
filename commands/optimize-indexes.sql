-- Optimización de índices para mejorar rendimiento
-- Ejecutar después de las migraciones principales

-- 1. Índices para tablas maestras (optimización de consultas IN)
-- Skill Categories
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_categories_id_active 
ON skill_categories (id) WHERE deleted_at IS NULL;

-- Skill Subcategories  
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_subcategories_id_active 
ON skill_subcategories (id) WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_subcategories_category_id 
ON skill_subcategories (category_id) WHERE deleted_at IS NULL;

-- Experience Levels
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_experience_levels_id_active 
ON experience_levels (id) WHERE deleted_at IS NULL;

-- Licenses (no tiene soft delete)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_licenses_id 
ON licenses (id);

-- 2. Índices para tablas de perfiles (optimización de consultas frecuentes)
-- Labour Profile Skills
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_labour_profile_skills_profile_id 
ON labour_profile_skills (labour_profile_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_labour_profile_skills_category_id 
ON labour_profile_skills (category_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_labour_profile_skills_subcategory_id 
ON labour_profile_skills (subcategory_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_labour_profile_skills_experience_id 
ON labour_profile_skills (experience_level_id);

-- User Licenses
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_licenses_user_id 
ON user_licenses (user_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_licenses_license_id 
ON user_licenses (license_id);

-- 3. Índices compuestos para consultas complejas
-- Labour profiles por usuario
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_labour_profiles_user_id 
ON labour_profiles (user_id) WHERE deleted_at IS NULL;

-- Builder profiles por usuario  
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_builder_profiles_user_id 
ON builder_profiles (user_id) WHERE deleted_at IS NULL;

-- 4. Índices para consultas de validación batch
-- Optimización específica para consultas WHERE id IN (...) AND deleted_at IS NULL
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_categories_batch_validation 
ON skill_categories (id) WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_subcategories_batch_validation 
ON skill_subcategories (id) WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_experience_levels_batch_validation 
ON experience_levels (id) WHERE deleted_at IS NULL;

-- 5. Índices adicionales para optimización extrema
-- Índices compuestos para consultas específicas
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_categories_id_deleted 
ON skill_categories (id, deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_skill_subcategories_id_deleted 
ON skill_subcategories (id, deleted_at) WHERE deleted_at IS NULL;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_experience_levels_id_deleted 
ON experience_levels (id, deleted_at) WHERE deleted_at IS NULL;

-- 5. Estadísticas de base de datos
ANALYZE skill_categories;
ANALYZE skill_subcategories;
ANALYZE experience_levels;
ANALYZE licenses;
ANALYZE labour_profile_skills;
ANALYZE user_licenses;
ANALYZE labour_profiles;
ANALYZE builder_profiles;
