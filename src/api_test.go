package main

import "context"

// fakeUserUseCase est un faux use case partage par les tests des handlers
// utilisateurs (api.users est type par l'interface userUseCase). Chaque test ne
// renseigne que la fonction dont il a besoin.
type fakeUserUseCase struct {
	registerFunc   func(ctx context.Context, pseudo, bio, ville string) (User, error)
	getProfileFunc func(ctx context.Context, id int) (User, error)
}

func (fake *fakeUserUseCase) Register(ctx context.Context, pseudo, bio, ville string) (User, error) {
	return fake.registerFunc(ctx, pseudo, bio, ville)
}

func (fake *fakeUserUseCase) GetProfile(ctx context.Context, id int) (User, error) {
	return fake.getProfileFunc(ctx, id)
}
