# PayPilot AI

Assistant intelligent pour freelances et petites entreprises : gestion de commandes, paiements PayPal et agent IA conversationnel capable d'agir sur vos données (avec confirmation humaine pour les actions sensibles).

## Fonctionnalités

- **Gestion des commandes** — création, consultation, mise à jour, suppression et détection des commandes en retard.
- **Paiements PayPal** — génération de liens de paiement, capture des paiements et réception des webhooks (sandbox par défaut).
- **Analytics** — résumé du chiffre d'affaires par période, montants en attente, évolution par rapport à la période précédente.
- **Agent IA (Claude)** — dialogue en langage naturel, avec des outils pour lister les commandes, consulter le détail d'une commande, calculer les revenus, rédiger des relances et créer des liens de paiement.
- **Human-in-the-loop** — les actions sensibles (création de lien de paiement, envoi de relance) sont mises en attente et doivent être confirmées ou rejetées par l'utilisateur.

## Stack

| Couche | Technologies |
|--------|--------------|
| Backend | Go 1.21, [chi](https://github.com/go-chi/chi), sqlx, SQLite |
| IA | API Anthropic (Claude) |
| Paiement | API REST PayPal |
| Frontend | React + TypeScript (à venir) |

## Démarrage rapide

### Prérequis

- Go 1.21+
- GCC (requis par `go-sqlite3`, qui utilise CGO)
- Un compte développeur PayPal (identifiants sandbox)
- Une clé API Anthropic

### Installation

```bash
cd backend
cp .env.example .env    # puis renseigner vos identifiants
make deps               # télécharger les dépendances
make run                # compiler et lancer le serveur
```

Le serveur démarre sur `http://localhost:8080`. La base SQLite et ses tables sont créées automatiquement au premier lancement.

### Commandes Make

| Commande | Description |
|----------|-------------|
| `make build` | Compile le binaire dans `bin/paypilot-api` |
| `make run` | Compile puis lance le serveur |
| `make dev` | Hot reload avec [air](https://github.com/air-verse/air) |
| `make dev-deps` | Installe air |
| `make test` | Lance les tests |
| `make env` | Crée `.env` à partir de `.env.example` |
| `make clean` | Supprime le binaire et la base locale |

## Configuration

Variables d'environnement (fichier `backend/.env`) :

| Variable | Description | Défaut |
|----------|-------------|--------|
| `PORT` | Port HTTP | `8080` |
| `ENV` | Environnement | `development` |
| `DATABASE_PATH` | Chemin de la base SQLite | `./paypilot.db` |
| `PAYPAL_CLIENT_ID` | Client ID PayPal | — |
| `PAYPAL_SECRET` | Secret PayPal | — |
| `PAYPAL_BASE_URL` | URL de l'API PayPal | `https://api-m.sandbox.paypal.com` |
| `ANTHROPIC_API_KEY` | Clé API Anthropic | — |
| `FRONTEND_URL` | Origine autorisée (CORS) et URL de retour PayPal | `http://localhost:5173` |

## API

| Méthode | Endpoint | Description |
|---------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/orders` | Lister les commandes |
| POST | `/api/orders` | Créer une commande |
| GET | `/api/orders/overdue` | Commandes en retard |
| GET | `/api/orders/{id}` | Détail d'une commande |
| PATCH | `/api/orders/{id}` | Mettre à jour une commande |
| DELETE | `/api/orders/{id}` | Supprimer une commande |
| GET | `/api/analytics/summary` | Résumé des revenus |
| GET | `/api/analytics/overdue` | Statistiques des impayés |
| POST | `/api/chat` | Envoyer un message à l'agent IA |
| GET | `/api/chat/history` | Historique de la conversation |
| DELETE | `/api/chat/history` | Effacer l'historique |
| GET | `/api/actions/pending` | Actions en attente de confirmation |
| POST | `/api/actions/{id}/confirm` | Confirmer une action |
| POST | `/api/actions/{id}/reject` | Rejeter une action |
| POST | `/api/paypal/create-order/{id}` | Créer une commande PayPal |
| POST | `/api/paypal/capture/{paypalOrderId}` | Capturer un paiement |
| POST | `/api/paypal/webhook` | Webhook PayPal |

### Exemples

Créer une commande :

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{"client_name":"Acme","client_email":"contact@acme.com","description":"Site vitrine","amount":1200,"currency":"EUR"}'
```

Discuter avec l'agent :

```bash
curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Quelles commandes sont encore en attente ce mois-ci ?"}'
```

Si l'agent propose une action sensible, la réponse contient un champ `pending_action` ; confirmez-la avec `POST /api/actions/{id}/confirm`.

## Architecture

```
backend/
├── cmd/server/main.go      # Point d'entrée, configuration du routeur
├── internal/
│   ├── config/             # Chargement de la configuration
│   ├── database/           # Connexion SQLite + migrations
│   ├── models/             # Structures de données
│   ├── repository/         # Requêtes base de données
│   ├── services/
│   │   ├── paypal/         # Client API PayPal
│   │   └── ai/             # Agent Claude + outils
│   ├── handlers/           # Handlers HTTP
│   └── middleware/         # Logger, recovery
└── pkg/response/           # Helpers de réponse JSON
```

## Auteur

Demetrius-ch
