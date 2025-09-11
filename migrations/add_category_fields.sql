-- Migration to add home description and key points fields to categories table
-- Run this SQL script to update your existing database

-- Add home screen description fields
ALTER TABLE categories ADD COLUMN home_description_en TEXT;
ALTER TABLE categories ADD COLUMN home_description_ar TEXT;

-- Add three key points fields
ALTER TABLE categories ADD COLUMN point1_en VARCHAR(255);
ALTER TABLE categories ADD COLUMN point1_ar VARCHAR(255);
ALTER TABLE categories ADD COLUMN point2_en VARCHAR(255);
ALTER TABLE categories ADD COLUMN point2_ar VARCHAR(255);
ALTER TABLE categories ADD COLUMN point3_en VARCHAR(255);
ALTER TABLE categories ADD COLUMN point3_ar VARCHAR(255);

-- Add image field
ALTER TABLE categories ADD COLUMN image_url VARCHAR(500);

-- Update existing categories with sample data (optional)
-- You can remove this section if you don't want sample data

UPDATE categories SET 
    home_description_en = 'Discover our amazing collection of ' || name_en || ' products',
    home_description_ar = 'اكتشف مجموعتنا المذهلة من منتجات ' || COALESCE(name_ar, name_en),
    point1_en = 'High quality products',
    point1_ar = 'منتجات عالية الجودة',
    point2_en = 'Competitive prices',
    point2_ar = 'أسعار تنافسية',
    point3_en = 'Fast delivery',
    point3_ar = 'توصيل سريع'
WHERE id > 0;
