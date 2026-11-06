-- جدول الوظائف (Jobs)
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NOT NULL,
    location VARCHAR(100) NOT NULL, -- مثل: Cairo, Remote, Hybrid
    salary_range VARCHAR(100),       -- مثل: $3,000 - $5,000 أو مصري
    description TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- جدول المتقدمين للوظائف (Applications)
CREATE TABLE applications (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    applicant_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    resume_url TEXT NOT NULL, -- رابط الـ CV أو الـ LinkedIn
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- قيد لمنع نفس الشخص يتقدم لنفس الوظيفة مرتين بنفس الإيميل
    CONSTRAINT unique_job_applicant UNIQUE (job_id, email)
);

-- Indexes لتحسين الأداء أثناء البحث والـ Joins
CREATE INDEX idx_jobs_location ON jobs(location);
CREATE INDEX idx_applications_job_id ON applications(job_id);