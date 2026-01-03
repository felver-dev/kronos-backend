# Guide de Configuration - ITSM Backend

## Prérequis

1. **Go 1.21+** installé
2. **MySQL 8.0+** installé et démarré
3. **Git** installé

## Installation

### 1. Installer les dépendances Go

```bash
go mod download
```

### 2. Configurer MySQL

#### Vérifier que MySQL est démarré

**Windows :**
```powershell
# Vérifier si MySQL est en cours d'exécution
Get-Service | Where-Object {$_.Name -like "*mysql*"}

# Démarrer MySQL si nécessaire
Start-Service MySQL80
# ou
net start MySQL80
```

**Linux/Mac :**
```bash
# Vérifier le statut
sudo systemctl status mysql
# ou
sudo service mysql status

# Démarrer MySQL si nécessaire
sudo systemctl start mysql
# ou
sudo service mysql start
```

#### Créer la base de données (optionnel - sera créée automatiquement)

```sql
CREATE DATABASE itsm_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. Configurer les variables d'environnement

Créez un fichier `.env` à la racine du dossier `backend` :

```bash
cp .env.example .env
```

Éditez le fichier `.env` avec vos paramètres MySQL :

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=votre_mot_de_passe
DB_NAME=itsm_db
```

**Note :** Si vous n'avez pas de mot de passe pour MySQL, laissez `DB_PASSWORD=` vide.

### 4. Exécuter les migrations

```bash
# Depuis le dossier backend
go run cmd/migrate/main.go

# Avec seeding (données initiales)
go run cmd/migrate/main.go -seed
```

### 5. Démarrer le serveur API

```bash
# Depuis le dossier backend
go run cmd/api/main.go
```

Le serveur sera accessible sur `http://localhost:8080`

## Vérification

### Vérifier que MySQL est accessible

**Windows (PowerShell) :**
```powershell
Test-NetConnection -ComputerName localhost -Port 3306
```

**Linux/Mac :**
```bash
mysql -u root -p -h localhost -P 3306
```

### Endpoints disponibles

- **API** : `http://localhost:8080/api/v1`
- **Health Check** : `http://localhost:8080/health`
- **Swagger UI** : `http://localhost:8080/swagger/index.html`

## Dépannage

### Erreur : "Aucune connexion n'a pu être établie"

**Solutions :**
1. Vérifiez que MySQL est démarré (voir section "Vérifier que MySQL est démarré")
2. Vérifiez les paramètres dans `.env` (host, port, user, password)
3. Vérifiez que le port 3306 n'est pas bloqué par un firewall
4. Essayez de vous connecter manuellement avec un client MySQL

### Erreur : "Access denied"

**Solutions :**
1. Vérifiez le nom d'utilisateur et le mot de passe dans `.env`
2. Vérifiez que l'utilisateur MySQL a les permissions nécessaires
3. Essayez de vous connecter manuellement avec les mêmes identifiants

### Erreur : "Unknown database"

La base de données sera créée automatiquement. Si l'erreur persiste, créez-la manuellement :

```sql
CREATE DATABASE itsm_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Support

Pour plus d'aide, consultez :
- Documentation MySQL : https://dev.mysql.com/doc/
- Documentation GORM : https://gorm.io/docs/

