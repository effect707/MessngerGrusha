package chat

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/effect707/MessngerGrusha/internal/domain"
)

type ChatRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chat, error)
	IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	CreateDirect(ctx context.Context, creatorID, recipientID uuid.UUID) (*domain.Chat, error)
	CreateGroup(ctx context.Context, name string, creatorID uuid.UUID, memberIDs []uuid.UUID) (*domain.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error)
	GetDirectChatID(ctx context.Context, userA, userB uuid.UUID) (uuid.UUID, error)
	AddMember(ctx context.Context, chatID, userID uuid.UUID, role domain.MemberRole) error
	RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error
	GetMembers(ctx context.Context, chatID uuid.UUID) ([]domain.ChatMember, error)
}

type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type UseCase struct {
	chatRepo ChatRepository
	userRepo UserRepository
}

func NewUseCase(chatRepo ChatRepository, userRepo UserRepository) *UseCase {
	return &UseCase{
		chatRepo: chatRepo,
		userRepo: userRepo,
	}
}

func (uc *UseCase) CreateDirectChat(ctx context.Context, creatorID, recipientID uuid.UUID) (*domain.Chat, error) {
	if creatorID == recipientID {
		return nil, fmt.Errorf("cannot create chat with yourself: %w", domain.ErrInvalidInput)
	}

	if _, err := uc.userRepo.GetByID(ctx, recipientID); err != nil {
		return nil, fmt.Errorf("recipient: %w", err)
	}

	existingChatID, err := uc.chatRepo.GetDirectChatID(ctx, creatorID, recipientID)
	if err == nil {
		var chat *domain.Chat
		chat, err = uc.chatRepo.GetByID(ctx, existingChatID)
		if err != nil {
			return nil, fmt.Errorf("get existing chat: %w", err)
		}
		return chat, nil
	}
	if err != domain.ErrNotFound {
		return nil, fmt.Errorf("check existing chat: %w", err)
	}

	chat, err := uc.chatRepo.CreateDirect(ctx, creatorID, recipientID)
	if err != nil {
		return nil, fmt.Errorf("create direct chat: %w", err)
	}
	return chat, nil
}

func (uc *UseCase) CreateGroupChat(ctx context.Context, name string, creatorID uuid.UUID, memberIDs []uuid.UUID) (*domain.Chat, error) {
	if name == "" {
		return nil, fmt.Errorf("group name: %w", domain.ErrInvalidInput)
	}

	for _, id := range memberIDs {
		if _, err := uc.userRepo.GetByID(ctx, id); err != nil {
			return nil, fmt.Errorf("member %s: %w", id, err)
		}
	}

	chat, err := uc.chatRepo.CreateGroup(ctx, name, creatorID, memberIDs)
	if err != nil {
		return nil, fmt.Errorf("create group chat: %w", err)
	}
	return chat, nil
}

func (uc *UseCase) GetChat(ctx context.Context, chatID, userID uuid.UUID) (*domain.Chat, error) {
	isMember, err := uc.chatRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, domain.ErrForbidden
	}

	chat, err := uc.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get chat: %w", err)
	}
	return chat, nil
}

func (uc *UseCase) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	chats, err := uc.chatRepo.GetUserChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user chats: %w", err)
	}

	for i := range chats {
		if chats[i].Type != domain.ChatTypeDirect || chats[i].Name != nil {
			continue
		}
		members, err := uc.chatRepo.GetMembers(ctx, chats[i].ID)
		if err != nil {
			continue
		}
		for _, m := range members {
			if m.UserID != userID {
				other, err := uc.userRepo.GetByID(ctx, m.UserID)
				if err == nil {
					name := other.DisplayName
					if name == "" {
						name = other.Username
					}
					chats[i].Name = &name
				}
				break
			}
		}
	}

	return chats, nil
}

func (uc *UseCase) AddMember(ctx context.Context, chatID, adderID, newMemberID uuid.UUID) error {
	chat, err := uc.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get chat: %w", err)
	}
	if chat.Type != domain.ChatTypeGroup {
		return fmt.Errorf("can only add members to group chats: %w", domain.ErrInvalidInput)
	}

	isMember, err := uc.chatRepo.IsMember(ctx, chatID, adderID)
	if err != nil {
		return fmt.Errorf("check adder membership: %w", err)
	}
	if !isMember {
		return domain.ErrForbidden
	}

	alreadyMember, err := uc.chatRepo.IsMember(ctx, chatID, newMemberID)
	if err != nil {
		return fmt.Errorf("check new member: %w", err)
	}
	if alreadyMember {
		return fmt.Errorf("user already a member: %w", domain.ErrAlreadyExists)
	}

	if _, err := uc.userRepo.GetByID(ctx, newMemberID); err != nil {
		return fmt.Errorf("new member: %w", err)
	}

	if err := uc.chatRepo.AddMember(ctx, chatID, newMemberID, domain.MemberRoleMember); err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	return nil
}

func (uc *UseCase) RemoveMember(ctx context.Context, chatID, removerID, targetID uuid.UUID) error {
	chat, err := uc.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("get chat: %w", err)
	}
	if chat.Type != domain.ChatTypeGroup {
		return fmt.Errorf("can only remove members from group chats: %w", domain.ErrInvalidInput)
	}

	isMember, err := uc.chatRepo.IsMember(ctx, chatID, removerID)
	if err != nil {
		return fmt.Errorf("check remover membership: %w", err)
	}
	if !isMember {
		return domain.ErrForbidden
	}

	if err := uc.chatRepo.RemoveMember(ctx, chatID, targetID); err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	return nil
}

func (uc *UseCase) GetMembers(ctx context.Context, chatID, userID uuid.UUID) ([]domain.ChatMember, error) {
	isMember, err := uc.chatRepo.IsMember(ctx, chatID, userID)
	if err != nil {
		return nil, fmt.Errorf("check membership: %w", err)
	}
	if !isMember {
		return nil, domain.ErrForbidden
	}

	members, err := uc.chatRepo.GetMembers(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}
	return members, nil
}
