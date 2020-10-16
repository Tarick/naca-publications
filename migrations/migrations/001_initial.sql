-- Write your migrate up statements here

{{ template "migrations/shared/trigger_set_timestamp.sql" . }}

CREATE TABLE publishers (
  uuid uuid primary key,
  name text unique not null,
  url text unique not null,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  modified_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE TRIGGER set_timestamp BEFORE UPDATE ON "publishers" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

create table publication_types (
  type text not null UNIQUE
  );

insert into publication_types (type) values('rss');
insert into publication_types (type) values('scrapped');
insert into publication_types (type) values('api');

create table world_languages (
  alpha2 varchar(2) primary key,
  name text not null UNIQUE,
  native_name text not null
);

insert into world_languages values
('ab', 'Abkhaz', 'аҧсуа'),
 ('aa', 'Afar', 'Afaraf'),
 ('af', 'Afrikaans', 'Afrikaans'),
 ('ak', 'Akan', 'Akan'),
 ('sq', 'Albanian', 'Shqip'),
 ('am', 'Amharic', 'አማርኛ'),
 ('ar', 'Arabic', 'العربية'),
 ('an', 'Aragonese', 'Aragonés'),
 ('hy', 'Armenian', 'Հայերեն'),
 ('as', 'Assamese', 'অসমীয়া'),
 ('av', 'Avaric', 'авар мацӀ, магӀарул мацӀ'),
 ('ae', 'Avestan', 'avesta'),
 ('ay', 'Aymara', 'aymar aru'),
 ('az', 'Azerbaijani', 'azərbaycan dili'),
 ('bm', 'Bambara', 'bamanankan'),
 ('ba', 'Bashkir', 'башҡорт теле'),
 ('eu', 'Basque', 'euskara, euskera'),
 ('be', 'Belarusian', 'Беларуская'),
 ('bn', 'Bengali', 'বাংলা'),
 ('bh', 'Bihari', 'भोजपुरी'),
 ('bi', 'Bislama', 'Bislama'),
 ('bs', 'Bosnian', 'bosanski jezik'),
 ('br', 'Breton', 'brezhoneg'),
 ('bg', 'Bulgarian', 'български език'),
 ('my', 'Burmese', 'ဗမာစာ'),
 ('ca', 'Catalan; Valencian', 'Català'),
 ('ch', 'Chamorro', 'Chamoru'),
 ('ce', 'Chechen', 'нохчийн мотт'),
 ('ny', 'Chichewa; Chewa; Nyanja', 'chiCheŵa, chinyanja'),
 ('zh', 'Chinese', '中文 (Zhōngwén), 汉语, 漢語'),
 ('cv', 'Chuvash', 'чӑваш чӗлхи'),
 ('kw', 'Cornish', 'Kernewek'),
 ('co', 'Corsican', 'corsu, lingua corsa'),
 ('cr', 'Cree', 'ᓀᐦᐃᔭᐍᐏᐣ'),
 ('hr', 'Croatian', 'hrvatski'),
 ('cs', 'Czech', 'česky, čeština'),
 ('da', 'Danish', 'dansk'),
 ('dv', 'Divehi; Dhivehi; Maldivian;', 'ދިވެހި'),
 ('nl', 'Dutch', 'Nederlands, Vlaams'),
 ('en', 'English', 'English'),
 ('eo', 'Esperanto', 'Esperanto'),
 ('et', 'Estonian', 'eesti, eesti keel'),
 ('ee', 'Ewe', 'Eʋegbe'),
 ('fo', 'Faroese', 'føroyskt'),
 ('fj', 'Fijian', 'vosa Vakaviti'),
 ('fi', 'Finnish', 'suomi, suomen kieli'),
 ('fr', 'French', 'français, langue française'),
 ('ff', 'Fula; Fulah; Pulaar; Pular', 'Fulfulde, Pulaar, Pular'),
 ('gl', 'Galician', 'Galego'),
 ('ka', 'Georgian', 'ქართული'),
 ('de', 'German', 'Deutsch'),
 ('el', 'Greek, Modern', 'Ελληνικά'),
 ('gn', 'Guaraní', 'Avañeẽ'),
 ('gu', 'Gujarati', 'ગુજરાતી'),
 ('ht', 'Haitian; Haitian Creole', 'Kreyòl ayisyen'),
 ('ha', 'Hausa', 'Hausa, هَوُسَ'),
 ('he', 'Hebrew', 'עברית'),
 ('hz', 'Herero', 'Otjiherero'),
 ('hi', 'Hindi', 'हिन्दी, हिंदी'),
 ('ho', 'Hiri Motu', 'Hiri Motu'),
 ('hu', 'Hungarian', 'Magyar'),
 ('ia', 'Interlingua', 'Interlingua'),
 ('id', 'Indonesian', 'Bahasa Indonesia'),
 ('ie', 'Interlingue', 'Originally called Occidental; then Interlingue after WWII'),
 ('ga', 'Irish', 'Gaeilge'),
 ('ig', 'Igbo', 'Asụsụ Igbo'),
 ('ik', 'Inupiaq', 'Iñupiaq, Iñupiatun'),
 ('io', 'Ido', 'Ido'),
 ('is', 'Icelandic', 'Íslenska'),
 ('it', 'Italian', 'Italiano'),
 ('iu', 'Inuktitut', 'ᐃᓄᒃᑎᑐᑦ'),
 ('ja', 'Japanese', '日本語 (にほんご／にっぽんご)'),
 ('jv', 'Javanese', 'basa Jawa'),
 ('kl', 'Kalaallisut, Greenlandic', 'kalaallisut, kalaallit oqaasii'),
 ('kn', 'Kannada', 'ಕನ್ನಡ'),
 ('kr', 'Kanuri', 'Kanuri'),
 ('ks', 'Kashmiri', 'कश्मीरी, كشميري\u200e'),
 ('kk', 'Kazakh', 'Қазақ тілі'),
 ('km', 'Khmer', 'ភាសាខ្មែរ'),
 ('ki', 'Kikuyu, Gikuyu', 'Gĩkũyũ'),
 ('rw', 'Kinyarwanda', 'Ikinyarwanda'),
 ('ky', 'Kirghiz, Kyrgyz', 'кыргыз тили'),
 ('kv', 'Komi', 'коми кыв'),
 ('kg', 'Kongo', 'KiKongo'),
 ('ko', 'Korean', '한국어 (韓國語), 조선말 (朝鮮語)'),
 ('ku', 'Kurdish', 'Kurdî, كوردی\u200e'),
 ('kj', 'Kwanyama, Kuanyama', 'Kuanyama'),
 ('la', 'Latin', 'latine, lingua latina'),
 ('lb', 'Luxembourgish, Letzeburgesch', 'Lëtzebuergesch'),
 ('lg', 'Luganda', 'Luganda'),
 ('li', 'Limburgish, Limburgan, Limburger', 'Limburgs'),
 ('ln', 'Lingala', 'Lingála'),
 ('lo', 'Lao', 'ພາສາລາວ'),
 ('lt', 'Lithuanian', 'lietuvių kalba'),
 ('lu', 'Luba-Katanga', ''),
 ('lv', 'Latvian', 'latviešu valoda'),
 ('gv', 'Manx', 'Gaelg, Gailck'),
 ('mk', 'Macedonian', 'македонски јазик'),
 ('mg', 'Malagasy', 'Malagasy fiteny'),
 ('ms', 'Malay', 'bahasa Melayu, بهاس ملايو\u200e'),
 ('ml', 'Malayalam', 'മലയാളം'),
 ('mt', 'Maltese', 'Malti'),
 ('mi', 'Māori', 'te reo Māori'),
 ('mr', 'Marathi (Marāṭhī)', 'मराठी'),
 ('mh', 'Marshallese', 'Kajin M̧ajeļ'),
 ('mn', 'Mongolian', 'монгол'),
 ('na', 'Nauru', 'Ekakairũ Naoero'),
 ('nv', 'Navajo, Navaho', 'Diné bizaad, Dinékʼehǰí'),
 ('nb', 'Norwegian Bokmål', 'Norsk bokmål'),
 ('nd', 'North Ndebele', 'isiNdebele'),
 ('ne', 'Nepali', 'नेपाली'),
 ('ng', 'Ndonga', 'Owambo'),
 ('nn', 'Norwegian Nynorsk', 'Norsk nynorsk'),
 ('no', 'Norwegian', 'Norsk'),
 ('ii', 'Nuosu', 'ꆈꌠ꒿ Nuosuhxop'),
 ('nr', 'South Ndebele', 'isiNdebele'),
 ('oc', 'Occitan', 'Occitan'),
 ('oj', 'Ojibwe, Ojibwa', 'ᐊᓂᔑᓈᐯᒧᐎᓐ'),
 ('cu', 'Old Church Slavonic, Church Slavic, Church Slavonic, Old Bulgarian, Old Slavonic', 'ѩзыкъ словѣньскъ'),
 ('om', 'Oromo', 'Afaan Oromoo'),
 ('or', 'Oriya', 'ଓଡ଼ିଆ'),
 ('os', 'Ossetian, Ossetic', 'ирон æвзаг'),
 ('pa', 'Panjabi, Punjabi', 'ਪੰਜਾਬੀ, پنجابی\u200e'),
 ('pi', 'Pāli', 'पाऴि'),
 ('fa', 'Persian', 'فارسی'),
 ('pl', 'Polish', 'polski'),
 ('ps', 'Pashto, Pushto', 'پښتو'),
 ('pt', 'Portuguese', 'Português'),
 ('qu', 'Quechua', 'Runa Simi, Kichwa'),
 ('rm', 'Romansh', 'rumantsch grischun'),
 ('rn', 'Kirundi', 'kiRundi'),
 ('ro', 'Romanian, Moldavian, Moldovan', 'română'),
 ('ru', 'Russian', 'русский язык'),
 ('sa', 'Sanskrit (Saṁskṛta)', 'संस्कृतम्'),
 ('sc', 'Sardinian', 'sardu'),
 ('sd', 'Sindhi', 'सिन्धी, سنڌي، سندھی\u200e'),
 ('se', 'Northern Sami', 'Davvisámegiella'),
 ('sm', 'Samoan', 'gagana faa Samoa'),
 ('sg', 'Sango', 'yângâ tî sängö'),
 ('sr', 'Serbian', 'српски језик'),
 ('gd', 'Scottish Gaelic; Gaelic', 'Gàidhlig'),
 ('sn', 'Shona', 'chiShona'),
 ('si', 'Sinhala, Sinhalese', 'සිංහල'),
 ('sk', 'Slovak', 'slovenčina'),
 ('sl', 'Slovene', 'slovenščina'),
 ('so', 'Somali', 'Soomaaliga, af Soomaali'),
 ('st', 'Southern Sotho', 'Sesotho'),
 ('es', 'Spanish; Castilian', 'español, castellano'),
 ('su', 'Sundanese', 'Basa Sunda'),
 ('sw', 'Swahili', 'Kiswahili'),
 ('ss', 'Swati', 'SiSwati'),
 ('sv', 'Swedish', 'svenska'),
 ('ta', 'Tamil', 'தமிழ்'),
 ('te', 'Telugu', 'తెలుగు'),
 ('tg', 'Tajik', 'тоҷикӣ, toğikī, تاجیکی\u200e'),
 ('th', 'Thai', 'ไทย'),
 ('ti', 'Tigrinya', 'ትግርኛ'),
 ('bo', 'Tibetan Standard, Tibetan, Central', 'བོད་ཡིག'),
 ('tk', 'Turkmen', 'Türkmen, Түркмен'),
 ('tl', 'Tagalog', 'Wikang Tagalog, ᜏᜒᜃᜅ᜔ ᜆᜄᜎᜓᜄ᜔'),
 ('tn', 'Tswana', 'Setswana'),
 ('to', 'Tonga (Tonga Islands)', 'faka Tonga'),
 ('tr', 'Turkish', 'Türkçe'),
 ('ts', 'Tsonga', 'Xitsonga'),
 ('tt', 'Tatar', 'татарча, tatarça, تاتارچا\u200e'),
 ('tw', 'Twi', 'Twi'),
 ('ty', 'Tahitian', 'Reo Tahiti'),
 ('ug', 'Uighur, Uyghur', 'Uyƣurqə, ئۇيغۇرچە\u200e'),
 ('uk', 'Ukrainian', 'українська'),
 ('ur', 'Urdu', 'اردو'),
 ('uz', 'Uzbek', 'zbek, Ўзбек, أۇزبېك\u200e'),
 ('ve', 'Venda', 'Tshivenḓa'),
 ('vi', 'Vietnamese', 'Tiếng Việt'),
 ('vo', 'Volapük', 'Volapük'),
 ('wa', 'Walloon', 'Walon'),
 ('cy', 'Welsh', 'Cymraeg'),
 ('wo', 'Wolof', 'Wollof'),
 ('fy', 'Western Frisian', 'Frysk'),
 ('xh', 'Xhosa', 'isiXhosa'),
 ('yi', 'Yiddish', 'ייִדיש'),
 ('yo', 'Yoruba', 'Yorùbá'),
 ('za', 'Zhuang, Chuang', 'Saɯ cueŋƅ, Saw cuengh');

create table publications (
  uuid uuid PRIMARY KEY,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  language_code varchar(2) NOT NULL REFERENCES world_languages(alpha2),
  publisher_uuid uuid NOT NULL REFERENCES publishers(uuid) ON DELETE CASCADE ON UPDATE CASCADE ,
  type text NOT NULL REFERENCES publication_types(type) ON DELETE CASCADE ON UPDATE CASCADE,
  created_at timestamptz NOT NULL DEFAULT NOW(),
  modified_at timestamptz NOT NULL DEFAULT NOW(),
  UNIQUE(name, publisher_uuid)
);

CREATE TRIGGER set_timestamp BEFORE UPDATE ON "publications" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

---- create above / drop below ----

DROP trigger set_timestamp ON "publications";

DROP trigger set_timestamp ON "publishers";

DROP FUNCTION trigger_set_timestamp;

DROP TABLE "publications";
DROP TABLE "world_languages"
DROP TABLE "publishers";

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
