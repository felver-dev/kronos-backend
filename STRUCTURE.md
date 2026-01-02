# Structure du Projet Backend

## Organisation des dossiers

```
backend/
├── cmd/
│   └── api/                    # Point d'entrée de l'application
│       └── main.go             # Fichier principal qui démarre le serveur
│
├── config/                     # Configuration de l'application
│   └── config.go               # Chargement des variables d'environnement
│
├── database/                    # Connexion à la base de données
│   └── connection.go           # Configuration GORM et connexion MySQL
│
├── internal/                    # Code interne de l'application
│   ├── models/                 # Modèles GORM (entités de la base de données)
│   ├── dto/                    # DTOs (Data Transfer Objects) pour les requêtes/réponses
│   ├── handlers/               # Handlers HTTP (contrôleurs)
│   ├── middleware/             # Middlewares (JWT, CORS, logging, etc.)
│   ├── services/               # Logique métier
│   ├── repositories/           # Accès aux données (couche d'abstraction)
│   └── utils/                  # Utilitaires (JWT, password, response, etc.)
│
├── migrations/                 # Migrations de base de données
│
├── go.mod                      # Dépendances Go
├── go.sum                      # Checksums des dépendances
├── .gitignore                  # Fichiers à ignorer par Git
├── .env.example                # Exemple de fichier d'environnement
├── README.md                   # Documentation principale
└── STRUCTURE.md                # Ce fichier
```

## Architecture

### Modèles (models/)
Contient toutes les structures GORM représentant les tables de la base de données :
- `user.go` - Utilisateurs
- `role.go` - Rôles
- `ticket.go` - Tickets
- `incident.go` - Incidents
- `time_entry.go` - Entrées de temps
- etc.

### DTOs (dto/)
Contient les Data Transfer Objects pour séparer les modèles de base de données des objets utilisés dans les requêtes/réponses HTTP :
- `auth_dto.go` - DTOs pour l'authentification
- `user_dto.go` - DTOs pour les utilisateurs
- `ticket_dto.go` - DTOs pour les tickets
- `timesheet_dto.go` - DTOs pour la gestion du temps
- `delay_dto.go` - DTOs pour les retards et justifications
- `response_dto.go` - DTOs pour les réponses génériques

### Handlers (handlers/)
Contient les handlers HTTP qui traitent les requêtes :
- `auth_handler.go` - Authentification
- `user_handler.go` - Gestion des utilisateurs
- `ticket_handler.go` - Gestion des tickets
- etc.

### Services (services/)
Contient la logique métier de l'application :
- `auth_service.go` - Logique d'authentification
- `ticket_service.go` - Logique de gestion des tickets
- `timesheet_service.go` - Logique de gestion du temps
- etc.

### Repositories (repositories/)
Couche d'abstraction pour l'accès aux données :
- `user_repository.go` - Accès aux utilisateurs
- `ticket_repository.go` - Accès aux tickets
- etc.

### Middleware (middleware/)
- `auth.go` - Authentification JWT
- `cors.go` - Configuration CORS
- `role.go` - Vérification des rôles (à créer)
- `logger.go` - Logging (à créer)

### Utils (utils/)
Fonctions utilitaires :
- `jwt.go` - Génération et validation de tokens JWT
- `password.go` - Hashage et vérification de mots de passe
- `response.go` - Formatage des réponses HTTP
- `pagination.go` - Gestion de la pagination (à créer)
- `file.go` - Gestion des fichiers uploads (à créer)

## Prochaines étapes

1. Créer les modèles GORM basés sur le schéma de base de données
2. Implémenter les handlers pour chaque module
3. Créer les services avec la logique métier
4. Implémenter les repositories pour l'accès aux données
5. Ajouter les middlewares de validation des rôles
6. Créer les utilitaires pour la gestion des fichiers
7. Implémenter tous les endpoints API selon la documentation

