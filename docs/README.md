# Documentation Swagger

Ce dossier contient la documentation Swagger générée automatiquement à partir des commentaires dans le code.

## Génération de la documentation

Pour régénérer la documentation Swagger après avoir modifié les commentaires dans les handlers :

```bash
swag init -g cmd/api/main.go -o ./docs
```

## Accès à l'interface Swagger

Une fois le serveur démarré, accédez à l'interface Swagger UI à l'adresse :

```
http://localhost:8080/swagger/index.html
```

## Utilisation

1. **Démarrer le serveur** :
   ```bash
   go run cmd/api/main.go
   ```

2. **Ouvrir Swagger UI** :
   - Ouvrez votre navigateur
   - Allez sur `http://localhost:8080/swagger/index.html`

3. **Tester les endpoints** :
   - Cliquez sur un endpoint pour voir les détails
   - Cliquez sur "Try it out"
   - Remplissez les paramètres si nécessaire
   - Cliquez sur "Execute" pour envoyer la requête
   - Consultez la réponse

## Authentification

Pour tester les endpoints protégés :

1. Connectez-vous via `/api/v1/auth/login` pour obtenir un token JWT
2. Cliquez sur le bouton "Authorize" en haut de la page Swagger
3. Entrez `Bearer <votre_token>` dans le champ (remplacez `<votre_token>` par le token obtenu)
4. Cliquez sur "Authorize"
5. Tous les endpoints protégés seront maintenant accessibles avec votre token

## Fichiers générés

- `docs.go` : Code Go généré contenant la documentation
- `swagger.json` : Documentation au format JSON
- `swagger.yaml` : Documentation au format YAML

Ces fichiers sont générés automatiquement et ne doivent pas être modifiés manuellement.

