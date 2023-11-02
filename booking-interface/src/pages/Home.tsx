import React from 'react';
import Footer from '../components/Footer';
import NavigationBar from '../components/Navigation';
import { Box, Button, Card, CardActions, CardContent, CardMedia, Grid, Link, Typography } from '@mui/material';


const Home: React.FC = () =>  {
    return (

        <div>
            <NavigationBar/>
            <Box
            sx={{
                backgroundImage: "url('https://images.unsplash.com/photo-1496711640096-993b5f5fa420?auto=format&fit=crop&q=80&w=1932&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D')",
                backgroundPosition: "center",
                backgroundSize: "cover",
                backgroundRepeat: "no-repeat",
                height: "100vh",
                width: "100%",
                display:'flex',
                flexDirection:'column',
                justifyContent:'center',
                alignItems: 'center',
            }}
>
    <Typography variant="h1" color={'white'}>New Era of booking flights</Typography>
<Button  LinkComponent={Link} href='/flights' variant="contained" sx={{maxWidth: '200px', maxHeight: '100px', minWidth: '200px', minHeight: '100px', marginBottom:-6, marginTop:2}}>Book now</Button>
</Box>
<Grid container spacing={12} sx={{paddingTop:10, paddingBottom:10}}>
    <Grid item xs={2} /> {/* Empty item with a width of 2 grid units */}
      <Grid item>
      <Card sx={{ maxWidth: 345 }}>
      <CardMedia
        sx={{ height: 140 }}
        image="https://images.unsplash.com/photo-1496711640096-993b5f5fa420?auto=format&fit=crop&q=80&w=1932&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        title="green iguana"
      />
      <CardContent>
        <Typography gutterBottom variant="h5" component="div">
          Lizard
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Lizards are a widespread group of squamate reptiles, with over 6,000
          species, ranging across all continents except Antarctica
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">Share</Button>
        <Button size="small">Learn More</Button>
      </CardActions>
    </Card>
      </Grid>
      <Grid item>
      <Card sx={{ maxWidth: 345 }}>
      <CardMedia
        sx={{ height: 140 }}
        image="https://images.unsplash.com/photo-1496711640096-993b5f5fa420?auto=format&fit=crop&q=80&w=1932&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        title="green iguana"
      />
      <CardContent>
        <Typography gutterBottom variant="h5" component="div">
          Lizard
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Lizards are a widespread group of squamate reptiles, with over 6,000
          species, ranging across all continents except Antarctica
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">Share</Button>
        <Button size="small">Learn More</Button>
      </CardActions>
    </Card>
      </Grid>
      <Grid item >
      <Card sx={{ maxWidth: 345 }}>
      <CardMedia
        sx={{ height: 140 }}
        image="https://images.unsplash.com/photo-1496711640096-993b5f5fa420?auto=format&fit=crop&q=80&w=1932&ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
        title="green iguana"
      />
      <CardContent>
        <Typography gutterBottom variant="h5" component="div">
          Lizard
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Lizards are a widespread group of squamate reptiles, with over 6,000
          species, ranging across all continents except Antarctica
        </Typography>
      </CardContent>
      <CardActions>
        <Button size="small">Share</Button>
        <Button size="small">Learn More</Button>
      </CardActions>
    </Card>
      </Grid>
    </Grid>
    <Footer/>
        </div>
    )

}
  export default Home