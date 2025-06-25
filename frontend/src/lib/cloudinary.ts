type UploadField = 'image_url' | 'product_variant.image_url';

export interface CloudinaryUploadResponse {
  secure_url: string;
  public_id: string;
  [key: string]: any;
}

export interface CloudinaryError {
  error?: { message: string };
}

export async function uploadImageWithSignature(
  file: File,
  field: UploadField,
  getSignature: () => Promise<{
    timestamp: string;
    signature: string;
    api_key: string;
    cloud_name: string;
  }>
): Promise<{ field: UploadField; url: string }> {
  const { timestamp, signature, api_key, cloud_name } = await getSignature();

  const formData = new FormData();
  formData.append('file', file);
  formData.append('api_key', api_key);
  formData.append('timestamp', timestamp);
  formData.append('signature', signature);
  formData.append('folder', 'products');

  const cloudinaryRes = await fetch(
    `https://api.cloudinary.com/v1_1/${cloud_name}/image/upload`,
    {
      method: 'POST',
      body: formData,
    }
  );

  const data: CloudinaryUploadResponse & CloudinaryError = await cloudinaryRes.json();
  if (!cloudinaryRes.ok) throw new Error(data.error?.message || 'Upload failed');

  return { field, url: data.secure_url };
}
