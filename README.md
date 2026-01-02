# ITSM Backend - API REST

Backend de l'application ITSM (IT Service Management) pour MCI CARE CI.

## Technologies

- **Go 1.21+**
- **Gin** - Framework web
- **GORM** - ORM pour Go
- **MySQL** - Base de données
- **JWT** - Authentification

## Structure du projet

```
backend/
├── cmd/
│   └── api/              # Point d'entrée de l'application
├── config/               # Configuration
├── database/             # Scripts de base de données
├── internal/
│   ├── models/           # Modèles GORM
│   ├── handlers/         # Handlers HTTP
│   ├── middleware/       # Middlewares (JWT, CORS, etc.)
│   ├── services/         # Logique métier
│   ├── repositories/     # Accès aux données
│   └── utils/            # Utilitaires
├── migrations/           # Migrations de base de données
└── go.mod                # Dépendances Go
```

## Installation

1. Installer les dépendances :
```bash
go mod download
```

2. Configurer les variables d'environnement :
```bash
cp .env.example .env
# Éditer .env avec vos paramètres
```

3. Créer la base de données MySQL :
```sql
CREATE DATABASE itsm_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

4. Lancer l'application :
```bash
go run cmd/api/main.go
```

## Endpoints API

Base URL : `http://localhost:8080/api/v1`

Voir `ENDPOINTS_API_COMPLETS.md` pour la documentation complète des endpoints.

## Rôles

- **DSI** : Accès total
- **RESPONSABLE_IT** : Supervision et validation
- **TECHNICIEN_IT** : Traitement des tickets

## Documentation

- Cahier des charges : `../CAHIER_DES_CHARGES_V1_MIS_A_JOUR.md`
- Endpoints API : `../ENDPOINTS_API_COMPLETS.md`
- Schéma BDD : `../SCHEMA_BASE_DE_DONNEES.md`

