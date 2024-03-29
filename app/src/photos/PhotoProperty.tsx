import React from 'react';
import { PhotoGroupData } from './model';
import { Box, Chip, Stack } from '@mui/material';

// 'default' | 'primary' | 'secondary' | 'error' | 'info' | 'success' | 'warning',
const PhotoProperty = ({photoGroup }: {photoGroup: PhotoGroupData}) => {
    return <Box>
        <div>{photoGroup.metaData.modelInfo}</div>
        <Stack direction="row" spacing={1}>
            {photoGroup.tags.map(tag => <Chip label={tag.name}  />)}
        </Stack>
    </Box>
}


export default PhotoProperty;