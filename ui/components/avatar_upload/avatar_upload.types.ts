export interface AvatarUploadProps {
  isOpen: boolean
  onOpenChange: (isOpen: boolean) => void
  initialUrl: string
  username: string
  onSave: (url: string) => Promise<void>
}
