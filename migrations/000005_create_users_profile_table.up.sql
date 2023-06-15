CREATE TABLE IF NOT EXISTS user_profiles (
     domain_user_id BIGINT NOT NULL PRIMARY KEY,
     first_name TEXT NOT NULL,
     last_name TEXT NOT NULL,
     phone_number TEXT,
     date_of_birth DATE NOT NULL,
     address TEXT,
     about_me TEXT,
     profile_pic_url TEXT DEFAULT 'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png',
     CONSTRAINT fk_user_profiles
         FOREIGN KEY (domain_user_id)
             REFERENCES user_domains(id)
             ON DELETE CASCADE
);
